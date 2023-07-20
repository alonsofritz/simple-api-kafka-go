package usecase

import "github.com/alonsofritz/simple-api-kafka-go/internal/entity"

type CreateProductInputDTO struct {
	Name  string  `json: "name"`
	Price float64 `json: "price"`
}

type CreateProductOutputDTO struct {
	ID    string
	Name  string
	Price float64
}

type CreateProductUseCase struct {
	// Passando a interface
	ProductRepository entity.ProductRepository
}

// Funcao construtora do tipo criado
func NewCreateProductUseCase(productRepository entity.ProductRepository) *CreateProductUseCase {
	return &CreateProductUseCase{ProductRepository: productRepository}
}

func (u *CreateProductUseCase) Execute(input CreateProductInputDTO) (*CreateProductOutputDTO, error) {
	product := entity.NewProduct(input.Name, input.Price)
	err := u.ProductRepository.Create(product)
	if err != nil {
		return nil, err
	}

	// Retorna os dados do tipo DTO de output, sem estar acoplado ao dominio, retornando apenas um objeto burro.
	return &CreateProductOutputDTO{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	}, nil
}
