package web

import (
	"encoding/json"
	"net/http"

	"github.com/alonsofritz/simple-api-kafka-go/internal/usecase"
)

type ProductHandlers struct {
	CreateProductUseCase *usecase.CreateProductUseCase
	ListProductsUseCase  *usecase.ListProductsUseCase
}

func NewProductHandlers(createProductUseCase *usecase.CreateProductUseCase, listProductUseCase *usecase.ListProductsUseCase) *ProductHandlers {
	return &ProductHandlers{
		CreateProductUseCase: createProductUseCase,
		ListProductsUseCase:  listProductUseCase,
	}
}

// Mesmo esquema que foi feito no kafka
// json.NewDecoder : pega os dados do body e decoda no input, ou seja,
// os dados que eu vou receber da requisicao http no json ele vai hidratar
// esse cara: CreateProductInputDTO, para recebermos como DTO

// Perceba que funciona como um controller, porem ainda nao sabe o que esta
// acontecendo internamente dentro da minha aplicação

// Perceba que esa implementacao funciona como uma camada para conversar com
// minha aplicação, dessa forma, seguindo o mesmo principio deve ser possivel
// fazer o mesmo para graphql, grpc, para receber xml. Dessa forma trabalhamos
// com "conectores" nao precisando necessariamente mexer internamente na
// aplicacao.
func (p *ProductHandlers) CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateProductInputDTO
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	output, err := p.CreateProductUseCase.Execute(input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	// posso tambem cuspir o json na tela como seguinte:
	json.NewEncoder(w).Encode(output)
}

func (p *ProductHandlers) ListProductsHandler(w http.ResponseWriter, r *http.Request) {
	output, err := p.ListProductsUseCase.Execute()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
