package routes

import (
	"taskly-backend/controller"
	"taskly-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(c *gin.Engine) {

	c.POST("/registro", controller.Registro)
	c.POST("/login", controller.Login)

	p := c.Group("/")
	p.Use(middlewares.ValidarToken)
	{
		p.POST("/equipos", controller.CrearEquipo)
		p.GET("/equipos", controller.ObtenerEquipos)
		p.GET("/equipos/:id", controller.ObtenerEquipo)
		p.DELETE("/equipos/:id", controller.EliminarEquipo)

		//miembros
		p.POST("/equipos/:id/miembros", controller.AgregarMiembro)
		p.GET("/equipos/:id/miembros", controller.ObtenerMiembros)
		p.DELETE("/equipos/:id/miembros/:usuarioId", controller.EliminarMiembro)

		//tareas

		p.POST("/equipos/:id/tareas", controller.CrearTarea)
		p.GET("/equipos/:id/tareas", controller.ObtenerTareas)
		p.GET("/tareas/:id", controller.ObtenerTarea)
		p.PUT("/tareas/:id", controller.EditarTarea)
		p.PATCH("/tareas/:id/estado", controller.CambiarEstado)
		p.DELETE("/tareas/:id", controller.EliminarTarea)

		//busqueda
		p.GET("/usuarios/buscar", controller.BuscarUsuario)

		//Actualizar perfil
		p.GET("/me", controller.ObtenerPerfil)
		p.PUT("/me", controller.ActualizarPerfil)

	}

}
