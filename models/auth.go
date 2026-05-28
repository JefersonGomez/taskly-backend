package models

// LoginRequest Datos para iniciar sesión
// @Description Credenciales requeridas para autenticación
type LoginRequest struct {
	Email    string `json:"email" example:"ana@example.com"`
	Password string `json:"password" example:"MiClave123"`
}

// LoginResponse Respuesta exitosa de login
// @Description Contiene el token JWT para autenticar requests posteriores
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// RegisterRequest Datos para registrar un nuevo usuario
// @Description Información requerida para crear una cuenta
type RegisterRequest struct {
	Nombre   string `json:"nombre" example:"Ana López"`
	Email    string `json:"email" example:"ana@example.com"`
	Password string `json:"password" example:"MiClave123"`
}

// UserProfileResponse Perfil público de usuario
// @Description Datos visibles del usuario autenticado
type UserProfileResponse struct {
	ID     uint   `json:"id" example:"1"`
	Nombre string `json:"nombre" example:"Ana López"`
	Email  string `json:"email" example:"ana@example.com"`
}

// UpdateProfileRequest Datos para actualizar perfil
// @Description Campos opcionales para actualizar. Solo los enviados serán modificados.
type UpdateProfileRequest struct {
	Nombre   string `json:"nombre,omitempty" example:"Ana López Actualizado"`
	Password string `json:"password,omitempty" example:"NuevaClave456"`
}

// UpdateEstadoRequest Datos para actualizar el estado de una tarea
// @Description Solo contiene el campo estado con valores permitidos
type UpdateEstadoRequest struct {
	Estado string `json:"estado" example:"en_progreso" enums:"backlog,pendiente,en_progreso,completada"`
}
