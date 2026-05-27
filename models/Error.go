package models

// ErrorResponse Respuesta estándar de error
type ErrorResponse struct {
	Error   string `json:"error" example:"Datos inválidos"`
	Message string `json:"message,omitempty" example:"Campo requerido"`
	Code    int    `json:"code,omitempty" example:"400"`
}

// MessageResponse Respuesta con mensaje de éxito
type MessageResponse struct {
	Message string `json:"message" example:"Operación exitosa"`
}
