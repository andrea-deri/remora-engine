package engine

import (
	"context"
	"fmt"
	"remora/pkg/config"
	"remora/pkg/coral/compiler"
	"remora/pkg/coral/profile"
	"remora/pkg/coral/runner"
	"remora/pkg/customerror"
	"remora/pkg/protocol/transfer/btp"
	"remora/pkg/protocol/transfer/htp"
	"strings"
	"sync/atomic"
	"time"

	"github.com/rs/zerolog/log"
)

// Engine defines the interface for the Tidal Engine, listing the methods
// that any concrete implementation must provide.
type Engine interface {

	// Init defines the interface of the API that permits to initialize the
	// Tidal Engine and all its internal data structures
	Init(configMap config.ConfigMap, benthosAPI btp.BenthosAPI)

	// Search looks up a Minnow in the search tree using HTTP request details
	// such as the URL path, method, and headers, and returns the search result.
	Search(request htp.SearchRequest) htp.SearchResult

	// Validate evaluates the input data against the behavior of the Minnow
	// identified by the given ID, returning the effect result and any error encountered.
	Validate(request htp.ValidationRequest) htp.ValidationResult
}

// TidalEngine is the core execution layer of the Tidal rule engine.
//
// It maintains internal data structures holding Minnow references in a search tree
// and a map of pre-compiled Minnow behavior profiles providing fast, read-only
// access during Search and Validate operations.
//
// The engine employs a double-buffered model:
//   - The *primary* block is used concurrently by all runtime operations.
//   - The *secondary* block is used for updates.
//
// When new data is loaded or modified, the engine rebuilds the secondary structures
// and atomically switches them with the primary ones. This ensures lock-free,
// race-free and low-latency access during computation.
//
// Search and Validate always operate on the current primary block and are
// never blocked by ongoing updates.
type TidalEngine struct {

	// primarySearchTree points to the primary search tree containing Minnow references.
	// This tree is used for reads during computation and is updated atomically by switching
	// with the secondary tree.
	primarySearchTree atomic.Pointer[SearchTree]

	// secondarySearchTree is the secondary search tree used to accumulate updates.
	// It is swapped with the primary tree when data are updated.
	secondarySearchTree SearchTree

	// primaryProfileMap points to the primary map of Minnow behavior profiles.
	// This map is read during computation and is updated atomically by switching
	// with the secondary map.
	primaryProfileMap atomic.Pointer[map[string]profile.Profile]

	// secondaryProfileMap holds the secondary map of Minnow behavior profiles.
	// It is swapped with the primary map when data are updated.
	secondaryProfileMap map[string]profile.Profile

	// lastUpdate stores the timestamp of the last update applied to the internal
	// data structures.
	lastUpdate time.Time

	// benthosAPI defines the interface for communication with the Benthos layer.
	// It is used to retrieve Minnow profiles, references and update events.
	benthosAPI btp.BenthosAPI
}

// Init initializes the TidalEngine by setting up the search tree, profile map
// and Benthos API client. It must be called before any Search or Validate operations.
//
// This method also starts a background goroutine to listen for Profile update events
// from the Benthos layer, using the configured polling interval.
func (engine *TidalEngine) Init(configMap config.ConfigMap, benthosAPI btp.BenthosAPI) {

	engine.benthosAPI = benthosAPI

	engine.initSearchTree()
	engine.initProfileMap()
	engine.lastUpdate = time.Now()

	// Launch a goroutine in order to listen for event on Profiles
	pollingTime := time.Duration(configMap.ReadInt("TIDAL_MINNOW_PROFILES_UPDATE_POLLING_SECONDS", 60)) * time.Second
	log.Info().Msgf("Set polling time to [%v] to listen for changes on Benthos Registry", pollingTime)
	go engine.listenForChanges(pollingTime)
}

// initSearchTree initializes the internal search tree structures used by the TidalEngine.
//
// Both primary and secondary trees are created and populated with Minnow references
// loaded from the Benthos API. After population, the secondary tree is atomically
// swapped with the primary tree to ensure the engine starts with a consistent,
// fully initialized search tree. A fresh secondary tree is retained for future updates.
func (engine *TidalEngine) initSearchTree() {

	// Initialize search tree with empty instances in order to avoid
	// nil references during switch operation on final part of init.
	engine.primarySearchTree.Store(&SearchTree{})
	engine.secondarySearchTree = SearchTree{}

	log.Trace().Msgf("Initialization of Minnow Reference Registry started...")
	newPrimarySearchTree := NewSearchTree()
	newSecondarySearchTree := NewSearchTree()

	// By default, read all Minnow references and construct the search trees.
	// In future, this will be changed in order to provide on-demand
	// Minnow reference loading.
	referenceTransferData := engine.benthosAPI.ReadAllReferences()
	for _, minnowReference := range referenceTransferData.References {
		newPrimarySearchTree.AddSearchPath(minnowReference.Path, minnowReference.Method, minnowReference.Id, minnowReference.SpecialHeaders)
		newSecondarySearchTree.AddSearchPath(minnowReference.Path, minnowReference.Method, minnowReference.Id, minnowReference.SpecialHeaders)
	}

	// The search tree initialization is finalized by setting created search tree
	// in engine's secondary reference. The switchSearchTrees method provide the
	// correct swap of this reference to engine's primary tree.
	engine.secondarySearchTree = newPrimarySearchTree
	engine.switchSearchTrees()
	engine.secondarySearchTree = newSecondarySearchTree
	log.Trace().Msgf("Initialization of the search tree in Tidal Engine completed! [%v]", newPrimarySearchTree)
}

// initProfileMap initializes the internal profile map structures used by the TidalEngine.
//
// Both primary and secondary maps are created and populated with Minnow profiles
// loaded from the Benthos API. After population, the secondary map is atomically
// swapped with the primary map to ensure the engine starts with a consistent,
// fully initialized profile map. A fresh secondary map is retained for future updates.
func (engine *TidalEngine) initProfileMap() {

	// Initialize profile map with empty instances in order to avoid
	// nil references during switch operation on final part of init.
	engine.primaryProfileMap.Store(&map[string]profile.Profile{})
	engine.secondaryProfileMap = make(map[string]profile.Profile)

	log.Trace().Msgf("Initialization of Minnow Profile Registry started...")

	// By default, read all Minnow profiles and construct the profile map.
	// In future, this will be changed in order to provide on-demand
	// Minnow profile loading.
	profileTransferData := engine.benthosAPI.ReadAllProfiles()
	newPrimaryProfileMap := make(map[string]profile.Profile)
	for _, profile := range profileTransferData.Profiles {
		compiledProfile := engine.compileProfile(profile)
		newPrimaryProfileMap[profile.Id] = compiledProfile
	}

	// The profile map initialization is finalized by setting created search tree
	// in engine's secondary reference. The switchProfileMaps method provide the
	// correct swap of this reference to engine's primary profile map.
	engine.secondaryProfileMap = newPrimaryProfileMap
	engine.switchProfileMaps()
	engine.secondaryProfileMap = profile.DeepCopy(engine.secondaryProfileMap)
	log.Trace().Msgf("Initialization of the Profile map in Tidal Engine completed!")
}

// Search retrieves a Minnow reference from the search tree based on HTTP request information,
// including URL path, method, and headers.
//
// Starting from the root node, the method traverses the search tree following the path
// segments defined by the request URL. Standard token child nodes are prioritized during
// traversal while wildcard nodes are considered if no exact match exists. Path parameters
// are extracted from wildcard nodes and included in the result.
//
// Once the terminal node is reached, it is validated against the request method and headers.
// If multiple resources exist under the same endpoint, the resource whose header constraints
// best match the request (e.g., SOAPAction or custom headers) is selected. If no suitable
// resource is found, the search result indicates failure.
//
// Algorithm overview:
//  1. Traverse the search tree from the root following the request path.
//  2. Prioritize standard token nodes; fall back to wildcard nodes when needed.
//  3. Validate the terminal node against the request headers.
//  4. Return a SearchResult indicating whether a valid resource was found and its details.
func (engine *TidalEngine) Search(request htp.SearchRequest) htp.SearchResult {

	pathParams := map[string]any{}
	primarySearchTree := engine.primarySearchTree.Load()
	currentNode := &primarySearchTree.Root

	tokenizedPath := strings.SplitSeq(strings.Trim(request.Path, "/"), "/")

	// Search for each extracted path token, until the child of the last token is found.
	// During the analysis, extract the path parameter if the node is a wildcard.
	// If no node is found for one of the analyzed tokens, return false.
	for token := range tokenizedPath {

		child := currentNode.FindInChildren(token)
		if child == nil {
			return htp.SearchResult{IsFound: false}
		}
		if child.Type == NODE_TOKEN_WILDCARD {
			pathParams[child.Token[1:]] = token
		}
		currentNode = child
	}

	// Validate terminal node by using input HTTP method.
	terminalNodeByHttpMethod := currentNode.FindInChildren(strings.ToLower(request.Method))
	if terminalNodeByHttpMethod == nil || terminalNodeByHttpMethod.Type != NODE_TOKEN_TERMINAL {
		return htp.SearchResult{IsFound: false}
	}

	// Resolve resource matching by required input header
	resource, isFound := terminalNodeByHttpMethod.GetResourceByMatchingHeader(request.Headers)
	if !isFound {
		return htp.SearchResult{IsFound: false}
	}

	return htp.SearchResult{IsFound: true, ResourceId: resource.ResourceId, PathParams: pathParams}
}

// Validate evaluates the input data against the behavior associated with the Minnow Profile
// identified by the provided ID and returns the resulting [profile.EffectResult] along
// with any error encountered.
//
// Steps performed by the method:
//  1. Retrieve the Minnow's behavior profile using the given ID from the internal profile map.
//  2. Evaluate the input data against the retrieved profile using the profile validation engine.
//     If evaluation fails, returns an EffectResult with IsValid=false and the error message.
//     If evaluation succeeds, casts the result to [profile.EffectResult] and sets IsValid=true.
//  3. Return the final [profile.EffectResult] and any error encountered during evaluation.
func (engine *TidalEngine) Validate(request htp.ValidationRequest) htp.ValidationResult {

	minnowId := request.MinnowId
	input := request.Input

	log.Debug().Msgf("Request arrived on Tidal Engine. ID: [%s] Input: [%v]", minnowId, input)

	foundProfile, _ := engine.getByID(minnowId)
	log.Trace().Msgf("Retrieved Minnow profile from ID: [%s] with Profile: [%v]", minnowId, foundProfile)

	// Evaluate the input data against the found Profile
	rawResult, _, err := runner.Evaluate(foundProfile, input)
	if err != nil {
		log.Fatal().Msgf("error: %s", err)
		return htp.ValidationResult{
			Result: profile.EffectResult{IsValid: false},
			Error:  err.Error(),
		}
	}

	if rawResult == nil {
		log.Debug().Msgf("Sent request from Tidal Engine: [Empty]")
		return htp.ValidationResult{
			Result: profile.EffectResult{IsValid: false},
		}
	}

	// Enrich the outcome and emit the response as EffectResult
	result := rawResult.(profile.EffectResult)
	result.IsValid = true

	log.Debug().Msgf("Sent request from Tidal Engine: [%v]", result)
	return htp.ValidationResult{
		Result: result,
	}
}

// listenForChanges continuously monitors the Benthos event collection for edits related
// to Minnow profiles and updates the in-memory data structures accordingly.
//
// The function runs as a polling routine with intervals defined by the input duration.
// During each iteration, it performs the following steps:
//  1. Retrieves all pending Minnow events from Benthos whose operation type matches
//     supported profile operations (add, remove, edit, activate).
//  2. Applies the necessary modifications to the secondary search tree and profile map,
//     such as adding, editing or removing Minnow references as required.
//  3. Atomically switches the primary and secondary blocks, ensuring that the updated
//     state becomes the active primary without blocking read operations.
func (engine *TidalEngine) listenForChanges(pollingTime time.Duration) {

	pollingTimerTicker := time.NewTicker(pollingTime)

	for timeEvent := range pollingTimerTicker.C {

		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		notEvaluatedEventsData := engine.benthosAPI.ReadAllEvents()
		notEvaluatedEvents := notEvaluatedEventsData.Events

		// Skip if no updates are present
		numberOfMinnowEvents := len(notEvaluatedEvents)
		if numberOfMinnowEvents == 0 {
			cancel()
			continue
		}

		log.Trace().Msgf("Detected [%v] Minnow events during polling event at time [%v]!", numberOfMinnowEvents, timeEvent)

		for _, eventToEvaluate := range notEvaluatedEvents {

			switch eventToEvaluate.Operation.Type {

			case string(btp.MINNOW_EVENT_TYPE_ADD):
				engine.addProfileInInternalState(eventToEvaluate, engine.secondaryProfileMap, &engine.secondarySearchTree)

			case string(btp.MINNOW_EVENT_TYPE_REMOVE):
				engine.removeProfileFromInternalState(eventToEvaluate, engine.secondaryProfileMap, &engine.secondarySearchTree)

			case string(btp.MINNOW_EVENT_TYPE_EDIT):
				engine.editProfileFromInternalState(eventToEvaluate, engine.secondaryProfileMap)

			case string(btp.MINNOW_EVENT_TYPE_PROFILE_ACTIVATION):
				engine.activateProfileInInternalState(eventToEvaluate)

			case string(btp.MINNOW_EVENT_TYPE_BEHAVIOR_ACTIVATION):
				engine.activateProfileBehaviorInInternalState(eventToEvaluate)

			default:
				log.Trace().Msgf("Invalid operation [%v] on Profile with ID [%v]", eventToEvaluate.Operation.Type, eventToEvaluate.MinnowId)
			}
		}

		// Update profile map and search tree by atomically switching structures
		engine.switchProfileMaps()
		engine.switchSearchTrees()

		cancel()
	}
}

// getByID retrieves a Minnow profile from the primary profile map using its identifier.
//
// It returns the corresponding profile if found. If no profile exists for the given ID,
// an error is returned.
func (engine *TidalEngine) getByID(id string) (profile.Profile, error) {

	profileMap := engine.primaryProfileMap.Load()
	value, isOk := (*profileMap)[id]
	if !isOk {
		return profile.Profile{}, customerror.NewError(customerror.ErrorTidalSearchProfileNotFound, id)
	}
	return value, nil
}

// addProfileInInternalState updates the Tidal Engine's internal state including a new Profile.
//
// It retrieves the Profile from Benthos layer, compiles it, stores it in the provided
// profile map and refreshes the associated path in the search tree. If the retrieval from Benthos layer
// fails or the profile is not uniquely returned, the update is aborted.
func (engine *TidalEngine) addProfileInInternalState(event btp.MinnowEvent, profileMap map[string]profile.Profile, searchTree *SearchTree) {

	log.Trace().Msgf("Detected [%v] operation on Profile with ID [%v]", event.Operation.Type, event.MinnowId)

	_, isFound := profileMap[event.MinnowId]
	if !isFound {
		log.Trace().Msgf("Profile with ID [%v] not found in profile map and it will be added", event.MinnowId)
	} else {
		log.Trace().Msgf("Profile with ID [%v] already found in profile map and it will be replaced", event.MinnowId)
	}

	profilesToEvaluateData := engine.benthosAPI.ReadSingleProfile(event.MinnowId)

	err := profilesToEvaluateData.Error
	if err == "" && profilesToEvaluateData.Size == 1 {

		// Retrieving the first, and only, raw profile and compiling it
		rawProfile := profilesToEvaluateData.Profiles[0]
		profileMap[event.MinnowId] = engine.compileProfile(rawProfile)

		// Delete the path in the search tree in order to re-generate it
		// without conflicts (i.e. duplication when special header is added,
		// change of HTTP method, etc)
		searchTree.RemoveSearchPath(event.Operation.Reference.Path, event.Operation.Reference.Method, event.Operation.Reference.SpecialHeaders)

		// Insert again the path in the search tree as fresh new
		searchTree.AddSearchPath(rawProfile.Path, rawProfile.Method, rawProfile.Id, rawProfile.SpecialHeaders)

	} else {

		log.Error().Msgf("An error occurred while retrieving profile with ID [%v]: %v", event.MinnowId, err)
	}
}

// removeProfileFromInternalState updates the Tidal Engine's internal state removing an existing Profile.
//
// If the profile does not exist in the provided profile map, the operation is skipped. Otherwise, the
// profile entry is deleted and the associated path in the search tree path is removed.
func (engine *TidalEngine) removeProfileFromInternalState(event btp.MinnowEvent, profileMap map[string]profile.Profile, searchTree *SearchTree) {

	log.Trace().Msgf("Detected [%v] operation on Profile with ID [%v]", event.Operation.Type, event.MinnowId)

	_, isFound := profileMap[event.MinnowId]
	if !isFound {
		log.Trace().Msgf("Profile with ID [%v] not found in profile map and it will not be removed", event.MinnowId)
		return
	}

	log.Trace().Msgf("Profile with ID [%v] found in profile map and it will be removed", event.MinnowId)
	delete(profileMap, event.MinnowId)

	// Delete the path in the search tree in order to completely remove it
	reference := event.Operation.Reference
	searchTree.RemoveSearchPath(reference.Path, reference.Method, reference.SpecialHeaders)

}

// editProfileFromInternalState updates the Tidal Engine's internal state updating an existing Profile.
//
// It retrieves the Profile from Benthos layer, re-compiles it and stores it in the provided
// profile map, replacing the old one. If the retrieval from Benthos layer fails or the profile is not
// uniquely returned, the update is aborted.
func (engine *TidalEngine) editProfileFromInternalState(event btp.MinnowEvent, profileMap map[string]profile.Profile) {

	log.Trace().Msgf("Detected [%v] operation on Profile with ID [%v]", event.Operation.Type, event.MinnowId)

	_, isFound := profileMap[event.MinnowId]
	if !isFound {
		log.Trace().Msgf("Profile with ID [%v] not found in profile map and it will be added", event.MinnowId)
	} else {
		log.Trace().Msgf("Profile with ID [%v] already found in profile map and it will be replaced", event.MinnowId)
	}

	profilesToEvaluateData := engine.benthosAPI.ReadSingleProfile(event.MinnowId)

	err := profilesToEvaluateData.Error
	if err == "" && profilesToEvaluateData.Size == 1 {
		// Retrieving the first, and only, raw profile and compiling it
		minnowProfile := engine.compileProfile(profilesToEvaluateData.Profiles[0])
		profileMap[event.MinnowId] = minnowProfile
	} else {
		log.Error().Msgf("An error occurred while retrieving profile with ID [%v]: %v", event.MinnowId, err)
	}
}

// TODO: currently not implemented
func (*TidalEngine) activateProfileInInternalState(event btp.MinnowEvent) {

	log.Trace().Msgf("Detected [%v] operation on Profile with ID [%v]", event.Operation.Type, event.MinnowId)
	log.Warn().Msgf("Edit [%v] currently not implemented!", event.Operation.Type)
}

// TODO: currently not implemented
func (*TidalEngine) activateProfileBehaviorInInternalState(event btp.MinnowEvent) {

	log.Trace().Msgf("Detected [%v] operation on Profile with ID [%v]", event.Operation.Type, event.MinnowId)
	log.Warn().Msgf("Edit [%v] currently not implemented!", event.Operation.Type)
}

// compileProfile converts a Minnow profile into an executable compiled Profile that can be evaluated
// by the validation engine.
//
// The compilation is performed by [compiler.Compile]. If this operation fails, the method panics and
// provide an error message that includes the profile ID.
func (engine *TidalEngine) compileProfile(profile btp.MinnowProfile) profile.Profile {

	compiledProfile, err := compiler.Compile(profile)
	if err != nil {
		panic(fmt.Sprintf("an error occurred while trying to compiling Profile with ID [%s]. %s", profile.Id, err))
	}
	return compiledProfile
}

// switchSearchTrees atomically updates the Engine's search trees using a double-buffer pattern.
//
// This method creates a deep copy of the secondary tree (containing all accumulated updates) and
// stores it as the new primary, ensuring that readers see a stable, consistent state.
// The last update timestamp is updated to track when the primary was switched.
func (engine *TidalEngine) switchSearchTrees() {

	copyOfUpdatedSearchTree := engine.secondarySearchTree.DeepCopy()
	engine.primarySearchTree.Store(&copyOfUpdatedSearchTree)
	engine.lastUpdate = time.Now()
}

// switchProfileMaps atomically updates the Engine's profile map using a double-buffer pattern.
//
// This method creates a deep copy of the secondary profile map (containing all accumulated updates)
// and stores it as the new primary, ensuring that readers see a stable, consistent state.
// The last update timestamp is updated to track when the primary was switched.
func (engine *TidalEngine) switchProfileMaps() {

	copyOfUpdatedProfileMap := profile.DeepCopy(engine.secondaryProfileMap)
	engine.primaryProfileMap.Store(&copyOfUpdatedProfileMap)
	engine.lastUpdate = time.Now()
}
