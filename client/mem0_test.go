package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bytectlgo/mem0-go/types"
	"github.com/stretchr/testify/assert"
)

func TestNewMemoryClient(t *testing.T) {
	// 测试无效的 API key
	_, err := NewMemoryClient(ClientOptions{})
	assert.Error(t, err)

	// 测试有效的客户端创建
	client, err := NewMemoryClient(ClientOptions{
		APIKey: "test-key",
	})
	assert.NoError(t, err)
	assert.NotNil(t, client)
}

func TestAddMemory(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/v1/memories/", r.URL.Path)
		assert.Equal(t, "Token test-key", r.Header.Get("Authorization"))

		// 返回测试响应
		response := []types.Memory{
			{
				ID:     "test-id",
				Memory: "test memory",
				UserID: "test-user",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// 创建客户端
	client, _ := NewMemoryClient(ClientOptions{
		APIKey: "test-key",
		Host:   server.URL,
	})

	// 测试添加内存
	_, err := client.Add("test memory", types.MemoryOptions{})
	assert.NoError(t, err)
}

func TestGetMemory(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/memories/test-id/", r.URL.Path)

		// 返回测试响应
		response := types.Memory{
			ID:     "test-id",
			Memory: "test memory",
			UserID: "test-user",
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// 创建客户端
	client, _ := NewMemoryClient(ClientOptions{
		APIKey: "test-key",
		Host:   server.URL,
	})

	// 测试获取内存
	memory, err := client.Get("test-id")
	assert.NoError(t, err)
	assert.Equal(t, "test-id", memory.ID)
	assert.Equal(t, "test memory", memory.Memory)
}

func TestSearchMemory(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/v1/memories/search/", r.URL.Path)
		assert.Equal(t, "test", r.URL.Query().Get("query"))

		// 返回测试响应
		response := []types.Memory{
			{
				ID:     "test-id",
				Memory: "test memory",
				UserID: "test-user",
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// 创建客户端
	client, _ := NewMemoryClient(ClientOptions{
		APIKey: "test-key",
		Host:   server.URL,
	})

	// 测试搜索内存
	results, err := client.Search("test", nil)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, "test-id", results[0].ID)
}

func TestDeleteMemory(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/v1/memories/test-id/", r.URL.Path)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// 创建客户端
	client, _ := NewMemoryClient(ClientOptions{
		APIKey: "test-key",
		Host:   server.URL,
	})

	// 测试删除内存
	err := client.Delete("test-id")
	assert.NoError(t, err)
}
