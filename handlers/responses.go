package handlers

// LoginResponse representa a resposta de sucesso do login.
// Ela inclui o token JWT e a data de expiração do token.
// A anotação `example` é usada para fornecer um exemplo de como a resposta deve ser
// estruturada na documentação da API.
// A data de expiração é formatada como uma string no formato RFC3339.
// Exemplo: "2025-06-27T10:50:43-03:00".
// // swagger:response LoginResponse
// A resposta de login bem-sucedido contém um token JWT e a data de expiração do token.
// A data de expiração é fornecida no formato RFC3339, que é o padrão para
// representar datas e horas em formato legível por humanos
// e é amplamente utilizado em APIs RESTful.
// A anotação `example` é usada para fornecer um exemplo de como a resposta deve ser
// estruturada na documentação da API.
type LoginResponse struct {
	Token  string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxLCJ0eXBlIjoiam9nYWRvciIsImV4cCI6MTcxOTQ5OTg0Mywib3JpZ19pYXQiOjE3MTkzMjcwNDN9.some_signature"`
	Expire string `json:"expire" example:"2025-06-27T10:50:43-03:00"`
}

// ErrorResponse representa uma resposta de erro genérica para a documentação.
// Ela inclui uma mensagem de erro que pode ser usada para informar o usuário
// sobre o que deu errado durante a operação solicitada.
// A anotação `example` é usada para fornecer um exemplo de como a resposta deve ser
// estruturada na documentação da API.
// swagger:response ErrorResponse
// A resposta de erro genérica contém uma mensagem de erro que pode ser usada para informar o
// usuário sobre o que deu errado durante a operação solicitada.
// A anotação `example` é usada para fornecer um exemplo de como a resposta deve ser
// estruturada na documentação da API.
// A mensagem de erro é uma string que descreve o erro ocorrido, como "Usuário não encontrado" ou "Credenciais inválidas".
// Ela deve ser clara e informativa para que o usuário possa entender o que deu errado e como corrigir o problema.
// Exemplo: "Usuário não encontrado" ou "Credenciais inválidas".
type ErrorResponse struct {
	Error string `json:"error" example:"Mensagem de erro"`
}

// SuccessResponse representa uma resposta de sucesso genérica com uma mensagem.
// Ela é usada para indicar que uma operação foi concluída com sucesso.
// A anotação `example` é usada para fornecer um exemplo de como a resposta deve ser
// estruturada na documentação da API.
// swagger:response SuccessResponse
// A resposta de sucesso genérica contém uma mensagem que indica que a operação foi concluída com
// sucesso. A mensagem é uma string que pode ser usada para informar o usuário sobre o resultado
// da operação solicitada, como "Operação realizada com sucesso" ou "Dados atualizados com sucesso".
// A anotação `example` é usada para fornecer um exemplo de como a resposta deve ser
// estruturada na documentação da API.
// A mensagem de sucesso é uma string que descreve o resultado da operação, como "Operação realizada com sucesso" ou "Dados atualizados com sucesso".
// Ela deve ser clara e informativa para que o usuário possa entender que a operação foi concluída com êxito.
// Exemplo: "Operação realizada com sucesso" ou "Dados atualizados com sucesso".
type SuccessResponse struct {
	Message string `json:"message" example:"Operação realizada com sucesso"`
}
