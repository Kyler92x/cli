package config

import (
	"bytes"
	"testing"
)

func TestLoadLocalConfig(t *testing.T) {
	path := "../../../examples/exampleconf/.tx/config"
	localCfg, err := loadLocalConfigFromPath(path)
	if err != nil {
		t.Error(err)
	}

	if localCfg.Resources[0].GetAPv3Id() != "o:__organization_slug__:p:__project_slug__:r:__resource_slug__" {
		t.Errorf(
			"APIv3 ID is wrong; got %+v, expected %+v",
			localCfg.Resources[0].GetAPv3Id(),
			"o:__organization_slug__:p:__project_slug__:r:__resource_slug__",
		)
	}

	expected := LocalConfig{
		Host: "https://app.transifex.com",
		Path: path,
		LanguageMappings: map[string]string{
			"de":    "de-Br",
			"pt_BR": "foo",
		},
		Resources: []Resource{{
			OrganizationSlug: "__organization_slug__",
			ProjectSlug:      "__project_slug__",
			ResourceSlug:     "__resource_slug__",
			FileFilter:       "locale/<lang>/ui.po",
			SourceFile:       "locale/ui.pot",
			SourceLanguage:   "en",
			Type:             "PO",
			LanguageMappings: map[string]string{
				"pt_PT": "pt-pt",
				"pt_BR": "pt-br",
			},
			Overrides: map[string]string{
				"pt-pt": "locale/other/pt_PT/ui.po",
				"fr_CA": "locale/other/fr_CA/ui.po",
			},
			MinimumPercentage: -1,
		}},
	}

	if !localConfigsEqual(localCfg, &expected) {
		t.Errorf(
			"Local config is wrong; got %+v, expected %+v",
			localCfg,
			expected,
		)
	}
}

func TestSaveAndLoadLocalConfig(t *testing.T) {
	expected := LocalConfig{
		Host: "My Host",
		LanguageMappings: map[string]string{
			"aa": "bb",
			"cc": "dd",
		},
		Resources: []Resource{
			{
				OrganizationSlug: "My Organization Slug",
				ProjectSlug:      "My Project Slug",
				ResourceSlug:     "My Resource Slug",
				FileFilter:       "My File Filter",
				SourceFile:       "My Source File",
				SourceLanguage:   "My Source Language",
				Type:             "My Type",
				LanguageMappings: map[string]string{
					"ee": "ff",
					"gg": "hh",
				},
				Overrides: map[string]string{
					"ee": "ff",
					"gg": "hh",
				},
			},
		},
	}

	var buffer bytes.Buffer
	err := expected.saveToWriter(&buffer)
	if err != nil {
		t.Error(err)
	}

	newLocalCfg, err := loadLocalConfigFromBytes(buffer.Bytes())
	if err != nil {
		t.Error(err)
	}

	if !localConfigsEqual(&expected, newLocalCfg) {
		t.Errorf(
			"Root config is wrong; got %+v, expected %+v",
			newLocalCfg,
			expected,
		)
	}
}

func TestChangeSaveAndLoadLocalConfig(t *testing.T) {
	initial := LocalConfig{
		Host: "My Host",
		Resources: []Resource{
			{
				OrganizationSlug: "My Organization Slug",
				ProjectSlug:      "My Project Slug",
				ResourceSlug:     "My Resource Slug",
				FileFilter:       "My File Filter",
				SourceFile:       "My Source File",
				SourceLanguage:   "My Source Language",
				Type:             "My Type",
				LanguageMappings: map[string]string{
					"ee": "ff",
					"gg": "hh",
				},
				Overrides: map[string]string{
					"ee": "ff",
					"gg": "hh",
				},
			},
		},
	}
	var buffer bytes.Buffer
	err := initial.saveToWriter(&buffer)
	if err != nil {
		t.Error(err)
	}

	// Load
	loaded, err := loadLocalConfigFromBytes(buffer.Bytes())
	if err != nil {
		t.Error(err)
	}

	// Change
	loaded.Resources[0].FileFilter = "My New File Filter"

	// Save again
	buffer.Reset()
	err = loaded.saveToWriter(&buffer)
	if err != nil {
		t.Error(err)
	}

	// Load again and check for file filter
	reloaded, err := loadLocalConfigFromBytes(buffer.Bytes())
	if err != nil {
		t.Error(err)
	}

	if reloaded.Resources[0].FileFilter != "My New File Filter" {
		t.Errorf(
			"Read wrong file_filter '%s', expected 'My New File Filter'",
			reloaded.Resources[0].FileFilter,
		)
	}

	loaded.Resources[0].MinimumPercentage = 10

	// Save again
	buffer.Reset()
	err = loaded.saveToWriter(&buffer)
	if err != nil {
		t.Error(err)
	}

	// Load again and check for file filter
	reloaded, err = loadLocalConfigFromBytes(buffer.Bytes())
	if err != nil {
		t.Error(err)
	}

	if reloaded.Resources[0].MinimumPercentage != 10 {
		t.Errorf(
			"Read wrong min_perc '%d', expected 10",
			reloaded.Resources[0].MinimumPercentage,
		)
	}
}
