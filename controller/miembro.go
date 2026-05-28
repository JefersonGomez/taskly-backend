package controller

import (
	"log"
	"strconv"
	"taskly-backend/models"
	"taskly-backend/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// AgregarMiembro godoc
// @Summary      Agregar un miembro a un equipo
// @Description  Crea un registro de miembro y lo asocia a un equipo existente. Solo el propietario del equipo puede realizar esta acción. Envía notificación por email al usuario agregado.
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
	// Obtener ID del equipo desde la ruta
	idEquipo := c.Param("id")
	equipoID, err := strconv.Atoi(idEquipo)
	if err != nil {
		c.JSON(400, gin.H{"error": "id incorrecto"})
		return
	}

	// Verificar que el equipo existe
	var equipo models.Equipo
	resultado := DB.First(&equipo, equipoID)
	if resultado.Error != nil {
		c.JSON(404, gin.H{"error": "equipo no encontrado"})
		return
	}

	// Bind del request body
	var miembro models.Miembro
	if err := c.BindJSON(&miembro); err != nil {
		c.JSON(400, gin.H{"error": "datos incorrectos"})
		return
	}

	// Validar rol
	if miembro.Rol != "admin" && miembro.Rol != "miembro" {
		c.JSON(400, gin.H{"error": "rol inválido, debe ser admin o miembro"})
		return
	}

	// Obtener datos del usuario que será agregado
	var usuarioAGregar models.Usuario
	if err := DB.First(&usuarioAGregar, miembro.UsuarioID).Error; err != nil {
		c.JSON(404, gin.H{"error": "usuario no encontrado"})
		return
	}

	// Obtener datos del usuario que está invitando (para el email)
	invitedByID := c.GetUint("id")
	var invitedByUser models.Usuario
	DB.First(&invitedByUser, invitedByID)
	invitedByName := invitedByUser.Nombre
	if invitedByName == "" {
		invitedByName = "Un miembro del equipo"
	}

	// Asignar equipo y crear miembro
	miembro.EquipoID = uint(equipoID)
	if err := DB.Create(&miembro).Error; err != nil {
		c.JSON(500, gin.H{"error": "error al agregar el miembro"})
		return
	}

	// 📧 ENVIAR EMAIL EN BACKGROUND (no bloquea la respuesta)
	go func() {
		// Pequeño delay para asegurar que la transacción de BD terminó
		time.Sleep(500 * time.Millisecond)

		err := utils.SendNotificationEmail(
			usuarioAGregar.Email,
			usuarioAGregar.Nombre,
			equipo.Nombre,
			miembro.Rol,
			invitedByName,
		)
		// Si falla el email, solo logueamos (no afectamos la respuesta principal)
		if err != nil {
			log.Printf("⚠️ No se pudo enviar email a %s: %v", usuarioAGregar.Email, err)
		}
	}()

	// Retornar respuesta con usuario preload
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
