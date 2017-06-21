// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * Copyright (C) 2016-2017 Canonical Ltd
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License version 3 as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 */

package service

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/CanonicalLtd/serial-vault/datastore"
	"github.com/CanonicalLtd/serial-vault/utils"
)

func TestReadConfig(t *testing.T) {
	config := ConfigSettings{}
	err := ReadConfig(&config, "../settings.yaml")
	if err != nil {
		t.Errorf("Error reading config file: %v", err)
	}
}

func TestReadConfigInvalidPath(t *testing.T) {
	config := ConfigSettings{}
	err := ReadConfig(&config, "not a good path")
	if err == nil {
		t.Error("Expected an error with an invalid config file.")
	}
}

func TestReadConfigInvalidFile(t *testing.T) {
	config := ConfigSettings{}
	err := ReadConfig(&config, "../README.md")
	if err == nil {
		t.Error("Expected an error with an invalid config file.")
	}
}

func TestFormatModelsResponse(t *testing.T) {
	var models []ModelSerialize
	models = append(models, ModelSerialize{ID: 1, BrandID: "Vendor", Name: "Alder 聖誕快樂"})
	models = append(models, ModelSerialize{ID: 2, BrandID: "Vendor", Name: "Ash"})

	w := httptest.NewRecorder()
	err := formatModelsResponse(true, "", "", "", models, w)
	if err != nil {
		t.Errorf("Error forming models response: %v", err)
	}

	var result ModelsResponse
	err = json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Error decoding the models response: %v", err)
	}
	if len(result.Models) != len(models) || !result.Success || result.ErrorMessage != "" {
		t.Errorf("Models response not as expected: %v", result)
	}
	if result.Models[0].Name != models[0].Name {
		t.Errorf("Expected the first model name of '%s', got: %s", models[0].Name, result.Models[0].Name)
	}
}

func TestFormatKeypairsResponse(t *testing.T) {
	var keypairs []datastore.Keypair
	keypairs = append(keypairs, datastore.Keypair{ID: 1, AuthorityID: "Vendor", KeyID: "12345678abcde", Active: true})
	keypairs = append(keypairs, datastore.Keypair{ID: 2, AuthorityID: "Vendor", KeyID: "abcdef123456", Active: true})

	w := httptest.NewRecorder()
	err := formatKeypairsResponse(true, "", "", "", keypairs, w)
	if err != nil {
		t.Errorf("Error forming keypairs response: %v", err)
	}

	var result KeypairsResponse
	err = json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Error decoding the keypairs response: %v", err)
	}
	if len(result.Keypairs) != len(keypairs) || !result.Success || result.ErrorMessage != "" {
		t.Errorf("Keypairs response not as expected: %v", result)
	}
	if result.Keypairs[0].KeyID != keypairs[0].KeyID {
		t.Errorf("Expected the first key ID '%s', got: %s", keypairs[0].KeyID, result.Keypairs[0].KeyID)
	}
}

func TestFormatSigningLogResponse(t *testing.T) {
	var signingLog []datastore.SigningLog
	signingLog = append(signingLog, datastore.SigningLog{ID: 1, Make: "System", Model: "Router 3400", SerialNumber: "A1", Fingerprint: "a1", Created: time.Now()})
	signingLog = append(signingLog, datastore.SigningLog{ID: 2, Make: "System", Model: "Router 3400", SerialNumber: "A2", Fingerprint: "a2", Created: time.Now()})

	w := httptest.NewRecorder()
	err := formatSigningLogResponse(true, "", "", "", signingLog, w)
	if err != nil {
		t.Errorf("Error forming signing log response: %v", err)
	}

	var result SigningLogResponse
	err = json.NewDecoder(w.Body).Decode(&result)
	if err != nil {
		t.Errorf("Error decoding the signing log response: %v", err)
	}
	if len(result.SigningLog) != len(signingLog) || !result.Success || result.ErrorMessage != "" {
		t.Errorf("Signing log response not as expected: %v", result)
	}
	if result.SigningLog[0].Fingerprint != signingLog[0].Fingerprint {
		t.Errorf("Expected the first fingerprint '%s', got: %s", signingLog[0].Fingerprint, result.SigningLog[0].Fingerprint)
	}
}

func TestRandomGeneration(t *testing.T) {
	n := 10
	// search for random strings enough smalls as to see if they are random
	tokens := make(map[string]string)
	for i := 0; i < n; i++ {
		// generate minimum amount of random data to verify it is enough random
		token, err := utils.GenerateRandomString(10)
		if err != nil {
			t.Errorf("Error generating random string: %v", err)
		}
		tokens[token] = token
	}

	// Check that we have n different tokens stored in the map
	if len(tokens) < n {
		t.Error("Generated random numbers are not unique")
	}
}

func TestModelName(t *testing.T) {
	name := "my-model-01"
	err := validateModelName(name)
	if err != nil {
		t.Errorf("Not a valid model name: %v", err)
	}
}

func TestModelNameEmpty(t *testing.T) {
	name := ""
	err := validateModelName(name)
	if err == nil {
		t.Error("Expected name not to be valid, but it is")
	}
	if err.Error() != "Name must not be empty" {
		t.Error("Error happening is not the one searched for")
	}
}

func TestModelNameInvalidChar(t *testing.T) {
	name := "my-model-01_"
	err := validateModelName(name)
	if err == nil {
		t.Error("Expected name not to be valid, but it is")
	}
	if !strings.Contains(err.Error(), "Name contains invalid characters") {
		t.Error("Error happening is not the one searched for")
	}
}

func TestModelNameUpperChar(t *testing.T) {
	name := "my-Model-01"
	err := validateModelName(name)
	if err == nil {
		t.Error("Expected name not to be valid, but it is")
	}
	if err.Error() != "Name must not contain uppercase characters" {
		t.Error("Error happening is not the one searched for")
	}
}
