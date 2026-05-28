package controller

import (
	"os"
	"taskly-backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Registro godoc
// @Summary      Registrar un nuevo usuario
// @Description  Crea una cuenta de usuario con email y contraseña. La contraseña se encripta automáticamente con bcrypt antes de guardar.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        usuario  body  models.RegisterRequest  true  "Datos de registro"
// @Success      200  {object}  models.Usuario  "Usuario registrado exitosamente"
// @Failure      400  {object}  models.ErrorResponse  "Datos inválidos o error al encriptar contraseña"
// @Failure      500  {object}  models.ErrorResponse  "Error interno del servidor"
// @Router       /registro [post]
func Registro(c *gin.Context) {
	var usuario models.Usuario
	if err := c.BindJSON(&usuario); err != nil {
		c.JSON(400, gin.H{"error": "Los datos son incorrectos"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(usuario.Password), 10)
	if err != nil {
		c.JSON(400, gin.H{"error": "No se pudo encriptar la contraseña"})
		return
	}

	usuario.Password = string(hash)
	DB.Create(&usuario)
	c.JSON(200, usuario)
}

// Login godoc
// @Summary      Iniciar sesión de usuario
// @Description  Autentica al usuario con email y contraseña. Retorna un token JWT válido por 24 horas para usar en endpoints protegidos.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials  body  models.LoginRequest  true  "Credenciales de acceso"
// @Success      200  {object}  models.LoginResponse  "Login exitoso"
// @Failure      400  {object}  models.ErrorResponse  "Datos de entrada inválidos"
// @Failure      404  {object}  models.ErrorResponse  "Usuario no encontrado o contraseña incorrecta"
// @Failure      500  {object}  models.ErrorResponse  "Error interno del servidor"
// @Router       /login [post]
func Login(c *gin.Context) {
	var body models.LoginRequest
	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Los datos son incorrectos"})
		return
	}

	var usuario models.Usuario
	DB.Where("email = ?", body.Email).First(&usuario)

	if usuario.ID == 0 {
		c.JSON(404, gin.H{"error": "No se encontró ninguna coincidencia"})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(body.Password))
	if err != nil {
		c.JSON(404, gin.H{"error": "No se encontraron coincidencias"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  usuario.ID,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		c.JSON(400, gin.H{"error": "No se pudo crear el token"})
		return
	}

	c.JSON(200, models.LoginResponse{Token: tokenString})
}

// BuscarUsuario godoc
// @Summary      Buscar usuario por email
// @Description  Permite buscar un usuario existente usando su email. Útil para verificar disponibilidad antes de registrar.
// @Tags         Usuarios
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        email  query  string  true  "Email del usuario a buscar"  example(ana@example.com)
// @Success      200  {object}  models.UserProfileResponse  "Usuario encontrado"
// @Failure      400  {object}  models.ErrorResponse  "Parámetro email faltante"
// @Failure      401  {object}  models.ErrorResponse  "No autenticado"
// @Failure      404  {object}  models.ErrorResponse  "Usuario no encontrado"
// @Failure      500  {object}  models.ErrorResponse  "Error interno del servidor"
// @Router       /usuarios/buscar [get]
func BuscarUsuario(c *gin.Context) {
	email := c.Query("email")
	var usuario models.Usuario
	resultado := DB.Where("email = ?", email).First(&usuario)

	// 🔧 CORRECCIÓN: GORM retorna error en resultado.Error
	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "no se encontró ninguna coincidencia con email"})
		return
	}

	c.JSON(200, models.UserProfileResponse{
		ID:     usuario.ID,
		Nombre: usuario.Nombre,
		Email:  usuario.Email,
	})
}

// ObtenerPerfil godoc
// @Summary      Obtener perfil del usuario autenticado
// @Description  Retorna los datos públicos del usuario que inició sesión, obtenidos desde el token JWT.
// @Tags         Usuarios
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  models.UserProfileResponse  "Perfil obtenido exitosamente"
// @Failure      401  {object}  models.ErrorResponse  "No autenticado"
// @Failure      404  {object}  models.ErrorResponse  "Usuario no encontrado"
// @Failure      500  {object}  models.ErrorResponse  "Error interno del servidor"
// @Router       /me [get]
func ObtenerPerfil(c *gin.Context) {
	idUsuario := c.GetUint("id")
	var usuario models.Usuario
	resultado := DB.First(&usuario, idUsuario)

	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "No se encontró el perfil de usuario logueado"})
		return
	}

	c.JSON(200, models.UserProfileResponse{
		ID:     usuario.ID,
		Nombre: usuario.Nombre,
		Email:  usuario.Email,
	})
}

// ActualizarPerfil godoc
// @Summary      Actualizar perfil del usuario autenticado
// @Description  Permite modificar el nombre y/o contraseña del usuario logueado. Solo los campos enviados serán actualizados.
// @Tags         Usuarios
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        datos  body  models.UpdateProfileRequest  false  "Campos a actualizar"
// @Success      200  {object}  models.UserProfileResponse  "Perfil actualizado exitosamente"
// @Failure      400  {object}  models.ErrorResponse  "Datos inválidos o error al encriptar"
// @Failure      401  {object}  models.ErrorResponse  "No autenticado"
// @Failure      404  {object}  models.ErrorResponse  "Usuario no encontrado"
// @Failure      500  {object}  models.ErrorResponse  "Error interno del servidor"
// @Router       /me [put]
func ActualizarPerfil(c *gin.Context) {
	idUsuario := c.GetUint("id")
	var body models.UpdateProfileRequest

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Los datos que ingresaste son incorrectos"})
		return
	}

	var usuario models.Usuario
	resultado := DB.First(&usuario, idUsuario)
	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "No se encontró el usuario"})
		return
	}

	if body.Nombre != "" {
		usuario.Nombre = body.Nombre
	}

	if body.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
		if err != nil {
			c.JSON(400, gin.H{"error": "No se pudo encriptar la contraseña"})
			return
		}
		usuario.Password = string(hash)
	}

	DB.Save(&usuario)

	c.JSON(200, models.UserProfileResponse{
		ID:     usuario.ID,
		Nombre: usuario.Nombre,
		Email:  usuario.Email,
	})
}
