package store

import (
	"io/ioutil"
	"os"
	"testing"
)

const ()

func TestLoadStore(t *testing.T) {
	yamlConf := []byte(`challenges:
  - name: challenge1
    id: id1
    image: image1
    memory: 100
  - name: challenge2
    id: id2
    image: image2
    memory: 200
`)

	type test struct {
		name      string
		predicate func(*testing.T, *storeDTO)
	}

	tests := []test{
		{
			name: "should load challenges",
			predicate: func(t *testing.T, dto *storeDTO) {
				if len(dto.Challenges) != 2 {
					t.Errorf("expected 2 challenges, got %d", len(dto.Challenges))
				}
			},
		},
		{
			name: "should load challenge1",
			predicate: func(t *testing.T, dto *storeDTO) {
				if dto.Challenges[0].Name != "challenge1" {
					t.Errorf("expected challenge1, got %s", dto.Challenges[0].Name)
				}
				if dto.Challenges[0].Id != "id1" {
					t.Errorf("expected id1, got %s", dto.Challenges[0].Id)
				}
				if dto.Challenges[0].Image != "image1" {
					t.Errorf("expected image1, got %s", dto.Challenges[0].Image)
				}
				if dto.Challenges[0].Memory != 100 {
					t.Errorf("expected 100, got %d", dto.Challenges[0].Memory)
				}
			},
		},
		{
			name: "should load challenge2",
			predicate: func(t *testing.T, dto *storeDTO) {
				if dto.Challenges[1].Name != "challenge2" {
					t.Errorf("expected challenge2, got %s", dto.Challenges[1].Name)
				}
				if dto.Challenges[1].Id != "id2" {
					t.Errorf("expected id2, got %s", dto.Challenges[1].Id)
				}
				if dto.Challenges[1].Image != "image2" {
					t.Errorf("expected image2, got %s", dto.Challenges[1].Image)
				}
				if dto.Challenges[1].Memory != 200 {
					t.Errorf("expected 200, got %d", dto.Challenges[1].Memory)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dto, err := loadStore(yamlConf)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			tt.predicate(t, dto)
		})
	}
}

func TestLoadStoreBadYaml(t *testing.T) {
	yamlConf := []byte(`challenges:
	  - name: challenge1
	    id: id1

	`)
	_, err := loadStore(yamlConf)
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestTransferChallenges(t *testing.T) {
	dto := &storeDTO{
		Challenges: []ChallengeTemplate{
			{
				Name:   "challenge1",
				Id:     "id1",
				Image:  "image1",
				Memory: 100,
			},
			{
				Name:   "challenge2",
				Id:     "id2",
				Image:  "image2",
				Memory: 200,
			},
		},
	}

	transferChallenges(dto)

	type test struct {
		name      string
		predicate func(*testing.T, *Store)
	}

	tests := []test{
		{
			name: "should transfer challenges",
			predicate: func(t *testing.T, s *Store) {
				if len(s.challenges) != 2 {
					t.Errorf("expected 2 challenges, got %d", len(s.challenges))
				}
			},
		},
		{
			name: "should transfer challenge1",
			predicate: func(t *testing.T, s *Store) {
				if s.challenges["challenge1"].Name != "challenge1" {
					t.Errorf("expected challenge1, got %s", s.challenges["challenge1"].Name)
				}
				if s.challenges["challenge1"].Id != "id1" {
					t.Errorf("expected id1, got %s", s.challenges["challenge1"].Id)
				}
				if s.challenges["challenge1"].Image != "image1" {
					t.Errorf("expected image1, got %s", s.challenges["challenge1"].Image)
				}
				if s.challenges["challenge1"].Memory != 100 {
					t.Errorf("expected 100, got %d", s.challenges["challenge1"].Memory)
				}
			},
		},
		{
			name: "should transfer challenge2",
			predicate: func(t *testing.T, s *Store) {
				if s.challenges["challenge2"].Name != "challenge2" {
					t.Errorf("expected challenge2, got %s", s.challenges["challenge2"].Name)
				}
				if s.challenges["challenge2"].Id != "id2" {
					t.Errorf("expected id2, got %s", s.challenges["challenge2"].Id)
				}
				if s.challenges["challenge2"].Image != "image2" {
					t.Errorf("expected image2, got %s", s.challenges["challenge2"].Image)
				}
				if s.challenges["challenge2"].Memory != 200 {
					t.Errorf("expected 200, got %d", s.challenges["challenge2"].Memory)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.predicate(t, store)
		})
	}
}

func TestInit(t *testing.T) {
	if store.challenges == nil {
		t.Error("challenges should not be nil")
	}
}

func TestGetChallenge(t *testing.T) {
	store = &Store{
		challenges: map[string]ChallengeTemplate{
			"challenge1": {
				Name:   "challenge1",
				Id:     "id1",
				Image:  "image1",
				Memory: 100,
			},
			"challenge2": {
				Name:   "challenge2",
				Id:     "id2",
				Image:  "image2",
				Memory: 200,
			},
		},
	}

	type test struct {
		name      string
		challenge string
		predicate func(*testing.T, *ChallengeTemplate)
	}

	tests := []test{
		{
			name:      "should get challenge1",
			challenge: "challenge1",
			predicate: func(t *testing.T, c *ChallengeTemplate) {
				if c.Name != "challenge1" {
					t.Errorf("expected challenge1, got %s", c.Name)
				}
				if c.Id != "id1" {
					t.Errorf("expected id1, got %s", c.Id)
				}
				if c.Image != "image1" {
					t.Errorf("expected image1, got %s", c.Image)
				}
				if c.Memory != 100 {
					t.Errorf("expected 100, got %d", c.Memory)
				}
			},
		},
		{
			name:      "should get challenge2",
			challenge: "challenge2",
			predicate: func(t *testing.T, c *ChallengeTemplate) {
				if c.Name != "challenge2" {
					t.Errorf("expected challenge2, got %s", c.Name)
				}
				if c.Id != "id2" {
					t.Errorf("expected id2, got %s", c.Id)
				}
				if c.Image != "image2" {
					t.Errorf("expected image2, got %s", c.Image)
				}
				if c.Memory != 200 {
					t.Errorf("expected 200, got %d", c.Memory)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			challenge, err := GetChallenge(tt.challenge)
			if err != nil {
				t.Errorf("unexpected error: %s", err)
			}

			tt.predicate(t, &challenge)
		})
	}
}

func TestChallengeExists(t *testing.T) {
	store = &Store{
		challenges: map[string]ChallengeTemplate{
			"challenge1": {
				Name:   "challenge1",
				Id:     "id1",
				Image:  "image1",
				Memory: 100,
			},
			"challenge2": {
				Name:   "challenge2",
				Id:     "id2",
				Image:  "image2",
				Memory: 200,
			},
		},
	}

	type test struct {
		name      string
		challenge string
		predicate func(*testing.T, bool)
	}

	tests := []test{
		{
			name:      "should return true for challenge1",
			challenge: "challenge1",
			predicate: func(t *testing.T, exists bool) {
				if !exists {
					t.Error("expected true, got false")
				}
			},
		},
		{
			name:      "should return true for challenge2",
			challenge: "challenge2",
			predicate: func(t *testing.T, exists bool) {
				if !exists {
					t.Error("expected true, got false")
				}
			},
		},
		{
			name:      "should return false for challenge3",
			challenge: "challenge3",
			predicate: func(t *testing.T, exists bool) {
				if exists {
					t.Error("expected false, got true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists := ChallengeExists(tt.challenge)
			tt.predicate(t, exists)
		})
	}
}

func TestLoad(t *testing.T) {
	// Write YAML config to file on disk
	data := []byte(`
challenges:
  - name: challenge1
    id: id1
    image: image1
    memory: 100
  - name: challenge2
    id: id2
    image: image2
    memory: 200
`)

	tmpfile, err := ioutil.TempFile("", "test_store_config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(data); err != nil {
		t.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	if err := Load(tmpfile.Name()); err != nil {
		t.Errorf("unexpected error: %s", err)
	}

	if store.challenges == nil {
		t.Error("challenges should not be nil")
	}

	if len(store.challenges) != 2 {
		t.Errorf("expected 2 challenges, got %d", len(store.challenges))
	}

	if store.challenges["challenge1"].Name != "challenge1" {
		t.Errorf("expected challenge1, got %s", store.challenges["challenge1"].Name)
	}

	if store.challenges["challenge1"].Id != "id1" {
		t.Errorf("expected id1, got %s", store.challenges["challenge1"].Id)
	}

	if store.challenges["challenge1"].Image != "image1" {
		t.Errorf("expected image1, got %s", store.challenges["challenge1"].Image)
	}

	if store.challenges["challenge1"].Memory != 100 {
		t.Errorf("expected 100, got %d", store.challenges["challenge1"].Memory)
	}

	if store.challenges["challenge2"].Name != "challenge2" {
		t.Errorf("expected challenge2, got %s", store.challenges["challenge2"].Name)
	}

	if store.challenges["challenge2"].Id != "id2" {
		t.Errorf("expected id2, got %s", store.challenges["challenge2"].Id)
	}

	if store.challenges["challenge2"].Image != "image2" {
		t.Errorf("expected image2, got %s", store.challenges["challenge2"].Image)
	}

	if store.challenges["challenge2"].Memory != 200 {
		t.Errorf("expected 200, got %d", store.challenges["challenge2"].Memory)
	}
}

func TestLoadInvalidConfig(t *testing.T) {
	// Write YAML config to file on disk
	data := []byte(`
challenges:
  - name: challenge1
`)
	tmpfile, err := ioutil.TempFile("", "test_store_config.yaml")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(data); err != nil {
		t.Fatal(err)
	}

	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	if err := Load(tmpfile.Name()); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestValidateStoreDTOWithNoNameChallenge(t *testing.T) {
	badTemplateNoName := ChallengeTemplate{
		Id:     "id1",
		Image:  "image1",
		Memory: 100,
	}

	dto := &storeDTO{
		Challenges: []ChallengeTemplate{badTemplateNoName},
	}

	if err := validateStoreDTO(dto); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestValidateStoreDTOWithNoIdChallenge(t *testing.T) {
	badTemplateNoId := ChallengeTemplate{
		Name:   "challenge1",
		Image:  "image1",
		Memory: 100,
	}

	dto := &storeDTO{
		Challenges: []ChallengeTemplate{badTemplateNoId},
	}

	if err := validateStoreDTO(dto); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestValidateStoreDTOWithNoImageChallenge(t *testing.T) {
	badTemplateNoImage := ChallengeTemplate{
		Name:   "challenge1",
		Id:     "id1",
		Memory: 100,
	}

	dto := &storeDTO{
		Challenges: []ChallengeTemplate{badTemplateNoImage},
	}

	if err := validateStoreDTO(dto); err == nil {
		t.Error("expected error, got nil")
	}
}

func TestValidateStoreDTOWithNoMemoryChallenge(t *testing.T) {
	badTemplateNoMemory := ChallengeTemplate{
		Name:  "challenge1",
		Id:    "id1",
		Image: "image1",
	}

	dto := &storeDTO{
		Challenges: []ChallengeTemplate{badTemplateNoMemory},
	}

	if err := validateStoreDTO(dto); err == nil {
		t.Error("expected error, got nil")
	}
}
