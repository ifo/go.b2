package b2

import (
	"encoding/json"
	"fmt"
	"testing"
)

func Test_Bucket_ListFileNames_Success(t *testing.T) {
	b := makeTestB2()
	bucket := makeTestBucket(b)

	fileAction := []Action{ActionUpload, ActionUpload, ActionUpload}
	setupFiles := ""
	for i := range fileAction {
		setupFiles += makeTestFileJson(i, fileAction[i])
		if i != len(fileAction)-1 {
			setupFiles += ","
		}
	}
	s := setupRequest(200, fmt.Sprintf(`{"files":[%s],"nextFileName":"name%d"}`, setupFiles, len(fileAction)))
	defer s.Close()

	response, err := bucket.ListFileNames("", 3)
	if err != nil {
		t.Fatalf("Expected no error, instead got %s", err)
	}

	if len(response.Files) != 3 {
		t.Fatalf("Expected two files, instead got %d", len(response.Files))
	}
	if response.NextFileName != fmt.Sprintf("name%d", len(fileAction)) {
		t.Errorf("Expected next file name to be name%d, instead got %s", len(fileAction), response.NextFileName)
	}
	if response.NextFileID != "" {
		t.Errorf("Expected no next file id, instead got %s", response.NextFileID)
	}
	for i, file := range response.Files {
		if file.Action != ActionUpload {
			t.Errorf("Expected action to be upload, instead got %v", file.Action)
		}
		if file.ID != fmt.Sprintf("id%d", i) {
			t.Errorf("Expected file ID to be id%d, instead got %s", i, fmt.Sprintf("id%d", i))
		}
		if file.Name != fmt.Sprintf("name%d", i) {
			t.Errorf("Expected file name to be name%d, instead got %s", i, fmt.Sprintf("name%d", i))
		}
		if file.Size != int64(10+i) {
			t.Errorf("Expected size to be %d, instead got %d", 10+i, file.Size)
		}
		if file.UploadTimestamp != int64(100+i) {
			t.Errorf("Expected upload timestamp to be %d, instead got %d", 10+i, file.UploadTimestamp)
		}
		if file.Bucket != bucket {
			t.Errorf("Expected file bucket to be bucket, instead got %+v", file.Bucket)
		}
	}
}

func Test_Bucket_ListFileNames_Errors(t *testing.T) {
	codes, bodies := errorResponses()
	b := makeTestB2()
	bucket := makeTestBucket(b)

	for i := range codes {
		s := setupRequest(codes[i], bodies[i])

		response, err := bucket.ListFileNames("", 0)
		testErrorResponse(err, codes[i], t)
		if response != nil {
			t.Errorf("Expected response to be empty, instead got %+v", response)
		}

		s.Close()
	}
}

func Test_Bucket_ListFileVersions_Success(t *testing.T) {
	b := makeTestB2()
	bucket := makeTestBucket(b)

	fileAction := []Action{ActionUpload, ActionHide, ActionStart}
	setupFiles := ""
	for i := range fileAction {
		setupFiles += makeTestFileJson(i, fileAction[i])
		if i != len(fileAction)-1 {
			setupFiles += ","
		}
	}
	s := setupRequest(200, fmt.Sprintf(`{"files":[%s],"nextFileId":"id%d","nextFileName":"name%d"}`,
		setupFiles, len(fileAction), len(fileAction)))
	defer s.Close()

	response, err := bucket.ListFileVersions("", "", 3)
	if err != nil {
		t.Fatalf("Expected no error, instead got %s", err)
	}

	if len(response.Files) != 3 {
		t.Fatalf("Expected three files, instead got %d", len(response.Files))
	}
	if response.NextFileName != "name3" {
		t.Errorf("Expected next file name to be name3, instead got %s", response.NextFileName)
	}
	if response.NextFileID != "id3" {
		t.Errorf("Expected next file id to be id3, instead got %s", response.NextFileID)
	}
	for i, file := range response.Files {
		if file.Action != fileAction[i] {
			t.Errorf("Expected action to be %v, instead got %v", fileAction[i], file.Action)
		}
		if file.ID != fmt.Sprintf("id%d", i) {
			t.Errorf("Expected file ID to be id%d, instead got %s", i, fmt.Sprintf("id%d", i))
		}
		if file.Name != fmt.Sprintf("name%d", i) {
			t.Errorf("Expected file name to be name%d, instead got %s", i, fmt.Sprintf("name%d", i))
		}
		if file.Size != int64(10+i) {
			t.Errorf("Expected size to be %d, instead got %d", 10+i, file.Size)
		}
		if file.UploadTimestamp != int64(100+i) {
			t.Errorf("Expected upload timestamp to be %d, instead got %d", 10+i, file.UploadTimestamp)
		}
		if file.Bucket != bucket {
			t.Errorf("Expected file bucket to be bucket, instead got %+v", file.Bucket)
		}
	}
}

func Test_Bucket_ListFileVersions_Errors(t *testing.T) {
	codes, bodies := errorResponses()
	b := makeTestB2()
	bucket := makeTestBucket(b)

	for i := range codes {
		s := setupRequest(codes[i], bodies[i])

		response, err := bucket.ListFileVersions("", "", 0)
		testErrorResponse(err, codes[i], t)
		if response != nil {
			t.Errorf("Expected response to be empty, instead got %+v", response)
		}

		s.Close()
	}
}

func makeTestFileJson(num int, action Action) string {
	file := FileMeta{
		ID:              fmt.Sprintf("id%d", num),
		Name:            fmt.Sprintf("name%d", num),
		Size:            int64(10 + num),
		ContentLength:   int64(10 + num),
		ContentSha1:     "sha1", // TODO make valid SHA1
		ContentType:     "text",
		Action:          action,
		FileInfo:        map[string]string{},
		UploadTimestamp: int64(100 + num),
	}
	fileJson, _ := json.Marshal(file)
	return string(fileJson)
}
