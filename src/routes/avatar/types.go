package avatar

import avatarv1 "github.com/mfdlabs/grid-service-websrv/avatar_v1"

// bodyColorsModel is the body colors model from avatar-fetch (old).
type bodyColorsModel struct {
	HeadColor     *int32 `json:"HeadColor"`
	TorsoColor    *int32 `json:"TorsoColor"`
	RightArmColor *int32 `json:"RightArmColor"`
	LeftArmColor  *int32 `json:"LeftArmColor"`
	RightLegColor *int32 `json:"RightLegColor"`
	LeftLegColor  *int32 `json:"LeftLegColor"`
}

// avatarFetchResponse is the response from the avatar fetch endpoint.
// This needs to be here due to the way that avatar-fetch's body colors
// model has changed.
type avatarFetchResponse struct {
	ResolvedAvatarType     *string                                             `json:"resolvedAvatarType,omitempty"`
	EquippedGearVersionIds []int64                                             `json:"equippedGearVersionIds,omitempty"`
	BackpackGearVersionIds []int64                                             `json:"backpackGearVersionIds,omitempty"`
	AssetAndAssetTypeIds   []avatarv1.RobloxApiAvatarModelsAssetIdAndTypeModel `json:"assetAndAssetTypeIds,omitempty"`
	AnimationAssetIds      *map[string]int64                                   `json:"animationAssetIds,omitempty"`
	BodyColors             *bodyColorsModel                                    `json:"bodyColors,omitempty"`
	Scales                 *avatarv1.RobloxWebResponsesAvatarScaleModel        `json:"scales,omitempty"`
}

// fromNewAvatarFetchResponse converts a new avatar fetch response to an old one.
func fromNewAvatarFetchResponse(newResponse *avatarv1.RobloxApiAvatarModelsAvatarFetchModel) *avatarFetchResponse {
	return &avatarFetchResponse{
		ResolvedAvatarType:     newResponse.ResolvedAvatarType,
		EquippedGearVersionIds: newResponse.EquippedGearVersionIds,
		BackpackGearVersionIds: newResponse.BackpackGearVersionIds,
		AssetAndAssetTypeIds:   newResponse.AssetAndAssetTypeIds,
		AnimationAssetIds:      newResponse.AnimationAssetIds,
		BodyColors: &bodyColorsModel{
			HeadColor:     newResponse.BodyColors.HeadColorId,
			TorsoColor:    newResponse.BodyColors.TorsoColorId,
			RightArmColor: newResponse.BodyColors.RightArmColorId,
			LeftArmColor:  newResponse.BodyColors.LeftArmColorId,
			RightLegColor: newResponse.BodyColors.RightLegColorId,
			LeftLegColor:  newResponse.BodyColors.LeftLegColorId,
		},
		Scales: newResponse.Scales,
	}
}
