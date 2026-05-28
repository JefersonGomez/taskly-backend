package controller

import (
	"strconv"
	"taskly-backend/models"

	"github.com/gin-gonic/gin"
)

// AgregarMiembro godoc
// @Summary      Agregar un miembro a un equipo
// @Description  Crea un registro de miembro y lo asocia a un equipo existente. Solo el propietario del equipo puede realizar esta acción.
// @Tags         Miembros
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path   int              true  "ID del equipo"          example(1)
// @Param        miembro  body   models.Miembro   true  "Datos del miembro a agregar"  example({"usuario_id": 5, "rol": "miembro"})
// @Success      201      {object}  models.Miembro  "Miembro agregado exitosamente"
// @Failure      400      {object}  models.ErrorResponse  "Datos inválidos o ID incorrecto"
// @Failure      401      {object}  models.ErrorResponse  "No autenticado"
// @Failure      404      {object}  models.ErrorResponse  "Equipo no encontrado"
// @Failure      500      {object}  models.ErrorResponse  "Error interno del servidor"
// @Router       /equipos/{id}/miembros [post]
func AgregarMiembro(c *gin.Context) {
	var equipo models.Equipo
	idEquipo := c.Param("id")
	id, err := strconv.Atoi(idEquipo)
	if err != nil {
		c.JSON(400, gin.H{"error": "id incorrecto"})
		return
	}
	resultado := DB.First(&equipo, id)
	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "equipo no encontrado"})
		return
	}
	var miembro models.Miembro
	if err := c.BindJSON(&miembro); err != nil {
		c.JSON(400, gin.H{"error": "datos incorrectos"})
		return
	}
	if miembro.Rol != "admin" && miembro.Rol != "miembro" {
		c.JSON(400, gin.H{"error": "rol inválido, debe ser admin o miembro"})
		return
	}
	miembro.EquipoID = equipo.ID
	DB.Create(&miembro)

	// retornar el miembro con el usuario cargado
	DB.Preload("Usuario").First(&miembro, miembro.ID)
	c.JSON(201, miembro)
}

// ObtenerMiembros godoc
// @Summary      Obtener miembros de un equipo
// @Description  Retorna la lista de usuarios que son miembros de un equipo específico, incluyendo su rol y datos básicos.
// @Tags         Miembros
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "ID del equipo"  example(1)
// @Success      200  {array}  models.Miembro  "Lista de miembros obtenida exitosamente"
// @Failure      400  {object}  models.ErrorResponse  "ID de equipo inválido"
// @Failure      401  {object}  models.ErrorResponse  "No autenticado"
// @Failure      404  {object}  models.ErrorResponse  "Equipo no encontrado"
// @Failure      500  {object}  models.ErrorResponse  "Error interno del servidor"
// @Router       /equipos/{id}/miembros [get]
func ObtenerMiembros(c *gin.Context) {
	idEquipo := c.Param("id")
	id, err := strconv.Atoi(idEquipo)
	if err != nil {
		c.JSON(400, gin.H{"error": "id incorrecto"})
		return
	}
	var equipo models.Equipo
	resultado := DB.First(&equipo, id)
	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "equipo no encontrado"})
		return
	}
	var miembros []models.Miembro
	DB.Preload("Usuario").Where("equipo_id = ?", id).Find(&miembros)
	c.JSON(200, miembros)
}

// EliminarMiembro godoc
// @Summary      Eliminar un miembro de un equipo
// @Description  Remueve permanentemente la relación de un usuario con un equipo. Solo el propietario del equipo puede realizar esta acción.
// @Tags         Miembros
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id         path  int  true  "ID del equipo"           example(1)
// @Param        usuarioId  path  int  true  "ID del usuario a eliminar"  example(5)
// @Success      200  {object}  models.MessageResponse  "Miembro eliminado exitosamente"
// @Failure      400  {object}  models.ErrorResponse  "IDs inválidos o formato incorrecto"
// @Failure      401  {object}  models.ErrorResponse  "No autenticado"
// @Failure      403  {object}  models.ErrorResponse  "Permisos insuficientes"
// @Failure      404  {object}  models.ErrorResponse  "Miembro no encontrado en este equipo"
// @Failure      500  {object}  models.ErrorResponse  "Error interno del servidor"
// @Router       /equipos/{id}/miembros/{usuarioId} [delete]
func EliminarMiembro(c *gin.Context) {
	var miembro models.Miembro
	idEquipo := c.Param("id")
	idUsuario := c.Param("usuarioId")
	id, err := strconv.Atoi(idEquipo)
	idUsu, err2 := strconv.Atoi(idUsuario)
	if err != nil || err2 != nil {
		c.JSON(400, gin.H{"error": "ids incorrectos"})
		return
	}
	resultado := DB.Where("equipo_id = ? AND usuario_id = ?", id, idUsu).First(&miembro)
	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "miembro no encontrado en ese equipo"})
		return
	}
	DB.Delete(&miembro)
	c.JSON(200, gin.H{"message": "miembro eliminado correctamente"})
}
