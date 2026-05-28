package controller

import (
	"strconv"
	"taskly-backend/models"
	"taskly-backend/utils"

	"github.com/gin-gonic/gin"
)

// CrearEquipo godoc
// @Summary      Crear un nuevo equipo de trabajo
// @Description  Crea un equipo asignando automáticamente al usuario autenticado (vía JWT) como propietario.
// @Tags         Equipos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        equipo  body      models.Equipo  true  "Datos para crear el equipo"
// @Success      201     {object}  models.Equipo  "Equipo creado exitosamente"
// @Failure      400     {object}  models.ErrorResponse  "Datos de entrada inválidos"
// @Failure      401     {object}  models.ErrorResponse  "No autenticado"
// @Failure      500     {object}  models.ErrorResponse  "Error interno del servidor"
// @Router       /equipos [post]
func CrearEquipo(c *gin.Context) {
	var equipo models.Equipo
	if err := c.BindJSON(&equipo); err != nil {
		c.JSON(400, models.ErrorResponse{Error: "datos incorrectos", Code: 400})
		return
	}
	equipo.OwnerID = c.GetUint("id")
	if err := DB.Create(&equipo).Error; err != nil {
		c.JSON(500, models.ErrorResponse{Error: "error al crear el equipo", Code: 500})
		return
	}
	c.JSON(201, equipo)
}

// ObtenerEquipos godoc
// @Summary      Listar equipos del usuario autenticado
// @Description  Retorna todos los equipos donde el usuario es propietario o miembro.
// @Tags         Equipos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   models.Equipo  "Lista de equipos obtenida exitosamente"
// @Failure      401  {object}  models.ErrorResponse  "No autenticado"
// @Failure      500  {object}  models.ErrorResponse  "Error interno del servidor"
// @Router       /equipos [get]
func ObtenerEquipos(c *gin.Context) {
	usuarioID := c.GetUint("id")

	var equiposDueno []models.Equipo
	if err := DB.Where("owner_id = ?", usuarioID).Find(&equiposDueno).Error; err != nil {
		c.JSON(500, models.ErrorResponse{Error: "error al consultar equipos", Code: 500})
		return
	}

	var miembros []models.Miembro
	if err := DB.Where("usuario_id = ?", usuarioID).Find(&miembros).Error; err != nil {
		c.JSON(500, models.ErrorResponse{Error: "error al consultar miembros", Code: 500})
		return
	}

	var equiposMiembro []models.Equipo
	if len(miembros) > 0 {
		equipoIDs := make([]uint, len(miembros))
		for i, m := range miembros {
			equipoIDs[i] = m.EquipoID
		}
		if err := DB.Where("id IN ?", equipoIDs).Find(&equiposMiembro).Error; err != nil {
			c.JSON(500, models.ErrorResponse{Error: "error al consultar equipos como miembro", Code: 500})
			return
		}
	}

	equiposMap := make(map[uint]models.Equipo)
	for _, e := range equiposDueno {
		equiposMap[e.ID] = e
	}
	for _, e := range equiposMiembro {
		equiposMap[e.ID] = e
	}

	var equipos []models.Equipo
	for _, e := range equiposMap {
		equipos = append(equipos, e)
	}

	c.JSON(200, equipos)
}

// ObtenerEquipo godoc
// @Summary      Obtener un equipo por ID
// @Description  Retorna los datos de un equipo específico, incluyendo miembros y tareas asociadas.
// @Tags         Equipos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "ID del equipo"
// @Success      200  {object}  models.Equipo  "Datos del equipo obtenidos exitosamente"
// @Failure      400  {object}  models.ErrorResponse  "ID inválido"
// @Failure      401  {object}  models.ErrorResponse  "No autenticado"
// @Failure      403  {object}  models.ErrorResponse  "Acceso denegado"
// @Failure      404  {object}  models.ErrorResponse  "Equipo no encontrado"
// @Failure      500  {object}  models.ErrorResponse  "Error interno del servidor"
// @Router       /equipos/{id} [get]
func ObtenerEquipo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, models.ErrorResponse{Error: "id incorrecto", Code: 400})
		return
	}

	var equipo models.Equipo
	resultado := DB.
		Preload("Miembros.Usuario").
		Preload("Tareas").
		First(&equipo, id)

	if resultado.Error != nil {
		c.JSON(404, models.ErrorResponse{Error: "equipo no encontrado", Code: 404})
		return
	}

	usuarioID := c.GetUint("id")
	if equipo.OwnerID != usuarioID {
		var miembro models.Miembro
		if err := DB.Where("equipo_id = ? AND usuario_id = ?", equipo.ID, usuarioID).First(&miembro).Error; err != nil {
			c.JSON(403, models.ErrorResponse{Error: "no tienes permisos para ver este equipo", Code: 403})
			return
		}
	}

	c.JSON(200, equipo)
}

// EliminarEquipo godoc
// @Summary      Eliminar un equipo por ID
// @Description  Elimina permanentemente un equipo. Solo el propietario (OwnerID) puede realizar esta acción.
// @Tags         Equipos
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id  path  int  true  "ID del equipo a eliminar"
// @Success      200  {object}  models.MessageResponse  "Equipo eliminado exitosamente"
// @Failure      400  {object}  models.ErrorResponse  "ID inválido"
// @Failure      401  {object}  models.ErrorResponse  "No autenticado"
// @Failure      403  {object}  models.ErrorResponse  "Permisos insuficientes"
// @Failure      404  {object}  models.ErrorResponse  "Equipo no encontrado"
// @Failure      500  {object}  models.ErrorResponse  "Error interno del servidor"
// @Router       /equipos/{id} [delete]
func EliminarEquipo(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(400, models.ErrorResponse{Error: "id incorrecto", Code: 400})
		return
	}

	var equipo models.Equipo
	resultado := DB.First(&equipo, id)
	if resultado.Error != nil {
		c.JSON(404, models.ErrorResponse{Error: "equipo no encontrado", Code: 404})
		return
	}

	if equipo.OwnerID != c.GetUint("id") {
		c.JSON(403, models.ErrorResponse{Error: "no tienes permisos para eliminar este equipo", Code: 403})
		return
	}

	if err := DB.Delete(&equipo).Error; err != nil {
		c.JSON(500, models.ErrorResponse{Error: "error al eliminar el equipo", Code: 500})
		return
	}

	c.JSON(200, models.MessageResponse{Message: "equipo eliminado correctamente"})
}

// TestEmail godoc
// @Summary Probar envío de email
// @Tags Debug
// @Router /debug/test-email [post]
func TestEmail(c *gin.Context) {
	err := utils.SendNotificationEmail(
		"gomesjeje504@gmail.com", // ← Tu email real
		"Usuario Test",
		"Equipo Test",
		"miembro",
		"Admin Taskly",
	)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Email enviado. Revisa tu bandeja"})
}
