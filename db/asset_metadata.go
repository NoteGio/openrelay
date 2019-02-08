package db

import (
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"github.com/notegio/openrelay/types"
	"github.com/jinzhu/gorm"
)

type AssetMetadata struct {
	AssetData   types.AssetData `gorm:"primary_key" json:"-"`
	RawMetadata string   `sql:"type:text;"  json:"-"`
	URI         string   `json:"token_uri"`
	Name        string	 `json:"name"`
	ExternalURL string	 `json:"external_ur,omitempty"`
	Image       string	 `json:"image,omitempty"`
	Description string   `sql:"type:text;" json:"description,omitempty"`
	BackgroundColor string `json:"background_color,omitempty"`
	Attributes  []AssetAttribute `json:"attributes,omitempty" gorm:"foreignkey:AssetData;association_foreigkey:AssetData"`
	RawAttributes string `json:"raw_attributes,omitempty"`
}

type AssetAttribute struct {
	AssetData   types.AssetData `gorm:"primary_key"  json:"-"`
	Name        string					`gorm:"primary_key"  json:"name"`
	Type        string					`json:"type,omitempty"`
	Value       string					`json:"value,omitempty"`
	DisplayType string 					`json:"display_type,omitempty"`
}

func toString(value interface{}) (string) {
	if value != nil {
		return value.(string)
	}
	return ""
}

func toAttributes(value interface{}) ([]AssetAttribute) {
	// Convert attributes from the myriad of formats they take for different
	// contracts
	results := []AssetAttribute{}
	switch v := value.(type) {
	case []interface{}:
		for _, item := range v {
			switch v2 := item.(type) {
			case map[string]interface{}:
				attr := AssetAttribute{}
				if val, ok := v2["trait_type"]; ok {
					attr.Name = fmt.Sprintf("%v", val)
				}
				if val, ok := v2["name"]; ok {
					attr.Name = fmt.Sprintf("%v", val)
				}
				if val, ok := v2["value"]; ok {
					switch v3 := val.(type) {
					case float64:
						attr.Type = "number"
						attr.Value = fmt.Sprintf("%f", v3)
					case string:
						attr.Type = "string"
						attr.Value = v3
					default:
						attr.Type = "object"
						encoded, err := json.Marshal(val)
						if err != nil {
							continue
						}
						attr.Value = string(encoded)
					}
				}
				if val, ok := v2["display_type"]; ok {
					attr.DisplayType = fmt.Sprintf("%v", val)
				}
				results = append(results, attr)
			default:
				log.Printf("Unknown attribute type: %T", v2)
			}
		}
	case map[string]interface{}:
		keys := make([]string, 0, len(v))
		for k := range v {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, key := range keys {
			item := v[key]
			switch v2 := item.(type) {
			case string:
				results = append(results, AssetAttribute{Name: key, Type: "string", Value: v2})
			case float64:
				results = append(results, AssetAttribute{Name: key, Type: "number", Value: fmt.Sprintf("%f", v2)})
			default:
				encoded, err := json.Marshal(v)
				if err == nil {
					results = append(results, AssetAttribute{Name: key, Type: "object", Value: string(encoded)})
				} // If there is an error, we're just going to skip it
			}
		}
	}
	return results
}

func (meta *AssetMetadata) UnmarshalJSON(data []byte) error {
	rawMetadata := make(map[string]interface{})
	if err := json.Unmarshal(data, &rawMetadata); err != nil {
		return err
	}
	meta.RawMetadata = string(data)
	meta.Name = toString(rawMetadata["name"])
	meta.ExternalURL = toString(rawMetadata["external_url"])
	meta.Image = toString(rawMetadata["image"])
	meta.Description = toString(rawMetadata["description"])
	meta.BackgroundColor = toString(rawMetadata["background_color"])
	meta.Attributes = toAttributes(rawMetadata["attributes"])
	encoding, err := json.Marshal(rawMetadata["attributes"])
	if err == nil {
		meta.RawAttributes = string(encoding)
	}
	return nil
}

// AfterFind populates Attribute Data upon querying for AssetMetadata. Gorm
// doesn't like our AssetData property as a primary key in Preload, as it tries
// to treat it as a collection of bytes instead of a key. This effectively
// implements preload implicitly, and accounts for Gorm's issues with byte
// slices and pointers.
func (meta *AssetMetadata) AfterFind(scope *gorm.Scope) {
	populateList := func(results *[]AssetMetadata) {

		// TODO: In theory, we ought to be able to get attributes using a subquery.
		// The code commented out here does that, but does not apply limits,
		// ordering, and offsets to the subquery, offsets to the subquery, so the
		// result set may be significantly larger than just what is related to these
		// items. For now, we're going to get the results and pass them back in,
		// rather than messing with subqueries. We may figure out later how to fix
		// offsets / limits

		// assetDataQuerySet := scope.DB().Debug().Table("asset_data").Select("asset_data")
		// if assetDataQuerySet.Error != nil {
		// 	fmt.Println(assetDataQuerySet.Error.Error())
		// }
		// db := scope.NewDB()
		// allAttributes := []AssetAttribute{}
		// if err := db.Debug().Model(&AssetAttribute{}).Where("asset_data IN (?)", assetDataQuerySet.QueryExpr()).Find(&allAttributes).Error; err != nil {
		// 	fmt.Println(err.Error())
		// }
		// assetDataMap := make(map[string][]AssetAttribute)
		// for _, attribute := range allAttributes {
		// 	assetDataMap[string(attribute.AssetData[:])] = append(assetDataMap[string(attribute.AssetData[:])], attribute)
		// }
		// for i, _ := range *results {
		// 	(*results)[i].Attributes = assetDataMap[string((*results)[i].AssetData)]
		// }
		assetDataSet := make(map[string]int)
		for i := range *results {
			assetDataSet[fmt.Sprintf("%#x", (*results)[i].AssetData[:])] = i
		}
		assetDataList := []*types.AssetData{}
		for _, value := range assetDataSet {
			assetDataList = append(assetDataList, &((*results)[value].AssetData))
		}
		db := scope.NewDB()
		allAttributes := []AssetAttribute{}
		db.Model(&AssetAttribute{}).Where("asset_data IN (?)", assetDataList).Find(&allAttributes)
		for _, attribute := range allAttributes {
			idx := assetDataSet[fmt.Sprintf("%#x", attribute.AssetData)]
			(*results)[idx].Attributes = append((*results)[idx].Attributes, attribute)
		}
	}
	switch results := scope.Value.(type) {
	case *[]AssetMetadata:
		populateList(results)
	case *AssetMetadata:
		resultList := &[]AssetMetadata{*results}
		populateList(resultList)
		*results = (*resultList)[0]
	}
}

func (meta *AssetMetadata) SetAssetData(assetData types.AssetData) {
	meta.AssetData = assetData
	for i := range meta.Attributes {
		meta.Attributes[i].AssetData = assetData
	}
}
