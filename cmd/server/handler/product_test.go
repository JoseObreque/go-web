package handler

import (
	"bytes"
	"encoding/json"
	"github.com/JoseObreque/go-web/internal/domain"
	"github.com/JoseObreque/go-web/internal/product"
	"github.com/JoseObreque/go-web/pkg/store"
	"github.com/JoseObreque/go-web/pkg/web"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func createServerForTestProducts(token string) *gin.Engine {
	// Token settings
	if token != "" {
		err := os.Setenv("TOKEN", token)
		if err != nil {
			panic(err)
		}
	}

	// Create a JSON store
	jsonStore := store.NewJsonStore("products_copy.json")

	// Obtains a slice of products
	products, err := jsonStore.GetAll()
	if err != nil {
		panic(err)
	}

	// Create a new product handler
	repository := product.NewRepository(products)
	service := product.NewService(repository)
	productHandler := NewProductHandler(service)

	// Define a new router
	router := gin.Default()

	// Add the product handler to the router
	productGroup := router.Group("/products")
	{
		productGroup.GET("/all", productHandler.GetAll())
		productGroup.GET("/:id", productHandler.GetById())
		productGroup.GET("/search", productHandler.GetByPriceGt())
		productGroup.POST("/new", productHandler.Create())
		productGroup.PUT("/:id", productHandler.FullUpdate())
		productGroup.PATCH("/:id", productHandler.PartialUpdate())
		productGroup.DELETE("/:id", productHandler.Delete())
	}

	return router
}

func createRequestTest(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	// Create a new request
	request := httptest.NewRequest(method, url, bytes.NewBuffer([]byte(body)))

	// Add elements to the header
	request.Header.Add("Content-Type", "application/json")

	return request, httptest.NewRecorder()
}

func TestProductHandler_GetAll_OK(t *testing.T) {
	router := createServerForTestProducts("")
	request, responseRecorder := createRequestTest(http.MethodGet, "https://localhost:8080/products/all", "")

	// Expected response
	jsonStore := store.NewJsonStore("products_copy.json")
	expectedResponse := web.Response{
		Data: []domain.Product{},
	}
	expectedProductsData, err := jsonStore.GetAll()
	if err != nil {
		panic(err)
	}
	expectedResponse.Data = expectedProductsData

	// Actual response
	router.ServeHTTP(responseRecorder, request)
	actualResponse := map[string][]domain.Product{}
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &actualResponse)

	// Assertions
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, expectedResponse.Data, actualResponse["data"])

}

func TestProductHandler_GetById_OK(t *testing.T) {
	router := createServerForTestProducts("")
	request, responseRecorder := createRequestTest(http.MethodGet, "https://localhost:8080/products/1", "")

	// Expected response
	jsonStore := store.NewJsonStore("products_copy.json")
	expectedResponse := web.Response{
		Data: domain.Product{},
	}
	expectedProductsData, err := jsonStore.GetOne(1)
	if err != nil {
		panic(err)
	}
	expectedResponse.Data = expectedProductsData

	// Actual response
	router.ServeHTTP(responseRecorder, request)
	actualResponse := map[string]domain.Product{}
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &actualResponse)

	// Assertions
	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, expectedResponse.Data, actualResponse["data"])

}

func TestProductHandler_Create_OK(t *testing.T) {
	// Expected response
	expectedResponse := web.Response{
		Data: domain.Product{
			Id:          501,
			Name:        "New Product",
			Quantity:    100,
			CodeValue:   "NewCode123",
			IsPublished: true,
			Expiration:  "25/10/2030",
			Price:       900,
		},
	}
	expectedProductData, err := json.Marshal(expectedResponse.Data)
	if err != nil {
		panic(err)
	}

	router := createServerForTestProducts("12345")
	request, responseRecorder := createRequestTest(
		http.MethodPost,
		"https://localhost:8080/products/new",
		string(expectedProductData),
	)
	request.Header.Add("token", "12345")

	// Actual response
	router.ServeHTTP(responseRecorder, request)
	actualResponse := map[string]domain.Product{}
	err = json.Unmarshal(responseRecorder.Body.Bytes(), &actualResponse)
	if err != nil {
		panic(err)
	}

	// Assertions
	assert.Equal(t, http.StatusCreated, responseRecorder.Code)
	assert.Equal(t, expectedResponse.Data, actualResponse["data"])
}

func TestProductHandler_Delete_OK(t *testing.T) {
	router := createServerForTestProducts("12345")
	request, responseRecorder := createRequestTest(
		http.MethodDelete,
		"https://localhost:8080/products/1",
		"",
	)
	request.Header.Add("token", "12345")

	// Actual response
	router.ServeHTTP(responseRecorder, request)
	actualResponse := responseRecorder.Body.Bytes()

	// Assertions
	assert.Equal(t, 204, responseRecorder.Code)
	assert.Nil(t, actualResponse)

}

func TestProductHandler_BadRequest(t *testing.T) {
	// Define a slice of http methods
	httpMethods := []string{
		http.MethodGet,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	}
	// Create a new router
	router := createServerForTestProducts("12345")

	// Iterate through the http methods slice
	for _, method := range httpMethods {
		// Create a request for every http method
		request, responseRecorder := createRequestTest(
			method,
			"https://localhost:8080/products/badId",
			"",
		)

		// Attach the token to the header and serve the request
		request.Header.Add("token", "12345")
		router.ServeHTTP(responseRecorder, request)

		// Unmarshal the actual response
		actualResponse := map[string]interface{}{}
		err := json.Unmarshal(responseRecorder.Body.Bytes(), &actualResponse)
		if err != nil {
			panic(err)
		}

		// Assertions
		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		assert.Equal(t, http.StatusText(http.StatusBadRequest), actualResponse["code"])
	}
}

func TestProductHandler_NotFound(t *testing.T) {
	// Define a slice of http methods
	httpMethods := []string{
		http.MethodGet,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
	}
	// Create a new router
	router := createServerForTestProducts("12345")

	// Create a body for the http methods that requires one
	newProduct := domain.Product{
		Name:        "New Product",
		Quantity:    100,
		CodeValue:   "NewCode123",
		IsPublished: true,
		Expiration:  "25/10/2030",
		Price:       900,
	}
	bodyProduct, err := json.Marshal(newProduct)
	if err != nil {
		panic(err)
	}

	// Iterate through the http methods slice
	for _, method := range httpMethods {
		// Create a request for every http method
		request, responseRecorder := createRequestTest(
			method,
			"https://localhost:8080/products/9999",
			string(bodyProduct),
		)

		// Attach the token to the header and serve the request
		request.Header.Add("token", "12345")
		router.ServeHTTP(responseRecorder, request)

		// Unmarshal the actual response
		actualResponse := map[string]interface{}{}
		err := json.Unmarshal(responseRecorder.Body.Bytes(), &actualResponse)
		if err != nil {
			panic(err)
		}

		// Assertions
		assert.Equal(t, http.StatusNotFound, responseRecorder.Code)
		assert.Equal(t, http.StatusText(http.StatusNotFound), actualResponse["code"])
	}
}

func TestProductHandler_Unauthorized(t *testing.T) {
	t.Run("Unauthorized PUT-PATCH-DELETE", func(t *testing.T) {
		// Define a slice of http methods
		httpMethods := []string{
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		}
		// Create a new router
		router := createServerForTestProducts("12345")

		// Create a body for the http methods that requires one
		newProduct := domain.Product{
			Name:        "New Product",
			Quantity:    100,
			CodeValue:   "NewCode123",
			IsPublished: true,
			Expiration:  "25/10/2030",
			Price:       900,
		}
		bodyProduct, err := json.Marshal(newProduct)
		if err != nil {
			panic(err)
		}

		// Iterate through the http methods slice
		for _, method := range httpMethods {
			// Create a request for every http method
			request, responseRecorder := createRequestTest(
				method,
				"https://localhost:8080/products/1",
				string(bodyProduct),
			)

			// Serve the request
			router.ServeHTTP(responseRecorder, request)

			// Unmarshal the actual response
			actualResponse := map[string]interface{}{}
			err := json.Unmarshal(responseRecorder.Body.Bytes(), &actualResponse)
			if err != nil {
				panic(err)
			}

			// Assertions
			assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
			assert.Equal(t, http.StatusText(http.StatusUnauthorized), actualResponse["code"])
		}
	})
	t.Run("Unauthorized POST", func(t *testing.T) {
		// Create a new router
		router := createServerForTestProducts("12345")

		// Create a body for the POST request
		newProduct := domain.Product{
			Name:        "New Product",
			Quantity:    100,
			CodeValue:   "NewCode123",
			IsPublished: true,
			Expiration:  "25/10/2030",
			Price:       900,
		}
		bodyProduct, err := json.Marshal(newProduct)
		if err != nil {
			panic(err)
		}

		// Create the POST request
		request, responseRecorder := createRequestTest(
			http.MethodPost,
			"https://localhost:8080/products/new",
			string(bodyProduct),
		)

		// Serve the request
		router.ServeHTTP(responseRecorder, request)

		// Unmarshal the actual response
		actualResponse := map[string]interface{}{}
		err = json.Unmarshal(responseRecorder.Body.Bytes(), &actualResponse)
		if err != nil {
			panic(err)
		}

		// Assertions
		assert.Equal(t, http.StatusUnauthorized, responseRecorder.Code)
		assert.Equal(t, http.StatusText(http.StatusUnauthorized), actualResponse["code"])
	})
}
