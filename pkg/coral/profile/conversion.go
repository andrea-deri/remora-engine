package profile

// DeepCopy returns an independent copy of a map[string]Profile. All nested structures
// (including slices and pointer fields) are duplicated, ensuring that updates to the
// copy do not affect the original map.
func DeepCopy(original map[string]Profile) map[string]Profile {

	if original == nil {
		return nil
	}

	copyMap := make(map[string]Profile, len(original))
	for key, profile := range original {
		copyMap[key] = deepCopyProfile(profile)
	}
	return copyMap
}

// deepCopyProfile creates a fully independent copy of a Profile instance. All contained behaviors
// are recreated to avoid sharing slice elements or pointer-based fields with the original.
func deepCopyProfile(profile Profile) Profile {

	copyBehaviors := make([]Behavior, len(profile.Behaviors))
	for i, behavior := range profile.Behaviors {
		copyBehaviors[i] = Behavior{
			Condition: deepCopyCondition(behavior.Condition),
			Effect:    deepCopyEffect(behavior.Effect),
		}
	}

	return Profile{Behaviors: copyBehaviors}
}

// deepCopyCondition produces a deep copy of a Condition value, recursively duplicating nested
// terms and re-allocating optional clause pointers. This ensures structural integrity when
// modifying derived condition trees.
func deepCopyCondition(condition Condition) Condition {

	copyCondition := Condition{Type: condition.Type}
	if condition.Clause != nil {
		copyClause := *condition.Clause
		copyCondition.Clause = &copyClause
	}

	if len(condition.Terms) > 0 {
		copyTerms := make([]Condition, len(condition.Terms))
		for i, term := range condition.Terms {
			copyTerms[i] = deepCopyCondition(term)
		}
		copyCondition.Terms = copyTerms
	}

	return copyCondition
}

// deepCopyEffect creates a deep copy of an Effect instance by duplicating its optional action
// and return pointers. The resulting value is safe to mutate without affecting the source.
func deepCopyEffect(effect Effect) Effect {

	copyEffect := Effect{}
	if effect.Action != nil {
		copyAction := *effect.Action
		copyEffect.Action = &copyAction
	}

	if effect.Return != nil {
		copyReturn := *effect.Return
		copyEffect.Return = &copyReturn
	}

	return copyEffect
}
