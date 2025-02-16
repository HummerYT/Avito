package integration

import (
	"bytes"
	"encoding/json"
	"github.com/google/uuid"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type RespRegister struct {
	Token string `json:"token"`
}

type userAuthIn struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,max=255"`
}

type requestSendCoin struct {
	ToUser string `json:"toUser" validate:"required"`
	Amount int64  `json:"amount" validate:"required,min=1"`
}

func TestBuyItem(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	urlPen := "http://localhost:8080/api/buy/pen"
	urlUmbrella := "http://localhost:8080/api/buy/umbrella"

	urlRegister := "http://localhost:8080/api/auth"
	userName1 := uuid.New().String()
	user := userAuthIn{
		Username: userName1,
		Password: "HardPass2007!",
	}
	data, err := json.Marshal(user)
	require.NoError(t, err)

	client := http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest("POST", urlRegister, bytes.NewReader(data))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var response RespRegister
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &response)
	require.NoError(t, err)

	req, err = http.NewRequest("GET", urlPen, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+response.Token)
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	req, err = http.NewRequest("GET", urlUmbrella, nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+response.Token)
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestGearUserCoin(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	urlGiftCoin := "http://localhost:8080/api/sendCoin"
	urlRegister := "http://localhost:8080/api/auth"

	userName1 := uuid.New().String()
	UserName2 := uuid.New().String()

	user1 := userAuthIn{
		Username: userName1,
		Password: "HardPass2007!",
	}
	data, err := json.Marshal(user1)
	require.NoError(t, err)

	client := http.Client{Timeout: time.Second * 5}
	req, err := http.NewRequest("POST", urlRegister, bytes.NewReader(data))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var response1 RespRegister
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &response1)
	require.NoError(t, err)

	user2 := userAuthIn{
		Username: UserName2,
		Password: "HardPass2007!",
	}
	data, err = json.Marshal(user2)
	require.NoError(t, err)

	client = http.Client{Timeout: time.Second * 5}
	req, err = http.NewRequest("POST", urlRegister, bytes.NewReader(data))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var response2 RespRegister
	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	err = json.Unmarshal(body, &response2)
	require.NoError(t, err)

	reqDataCoin := requestSendCoin{
		ToUser: UserName2,
		Amount: 1000,
	}
	data, err = json.Marshal(reqDataCoin)
	require.NoError(t, err)

	req, err = http.NewRequest("POST", urlGiftCoin, bytes.NewReader(data))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+response1.Token)
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)

	reqDataCoin = requestSendCoin{
		ToUser: UserName2,
		Amount: 1000,
	}
	data, err = json.Marshal(reqDataCoin)
	require.NoError(t, err)

	req, err = http.NewRequest("POST", urlGiftCoin, bytes.NewReader(data))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+response1.Token)
	resp, err = client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
