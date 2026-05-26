package controller

import (
	"os"
	"taskly-backend/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

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

func Login(c *gin.Context) {

	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BindJSON(&body); err != nil {
		c.JSON(400, gin.H{"error": "Los datos son incorrectos"})
		return
	}

	var usuario models.Usuario

	DB.Where("email = ?", body.Email).First(&usuario)

	if usuario.ID == 0 {
		c.JSON(404, gin.H{"error": "No se encontro ninguna concidencia"})
		return

	}

	err := bcrypt.CompareHashAndPassword([]byte(usuario.Password), []byte(body.Password))

	if err != nil {
		c.JSON(404, gin.H{"error": "No se encontraron concidencias"})
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

	c.JSON(200, gin.H{"token": tokenString})

}

func BuscarUsuario(c *gin.Context) {
	email := c.Query("email")

	var usuario models.Usuario

	resultado := DB.Where("email = ?", email).First(&usuario)

	if resultado != nil {
		c.JSON(404, gin.H{"error": "no se encontro ninguna concidencia con email"})
		return
	}

	c.JSON(200, gin.H{
		"Id":     usuario.ID,
		"nombre": usuario.Nombre,
		"email":  usuario.Email,
	})

}
