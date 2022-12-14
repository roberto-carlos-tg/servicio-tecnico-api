package controllers

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/roberto-carlos-tg/servicio-tecnico-api/database"
	"github.com/roberto-carlos-tg/servicio-tecnico-api/models"
)

type DetalleOrdenServicio struct {
	ID                  int                 `json:"id"`
	OrdenServicio       OrdenServicio       `json:"orden_servicio"`
	Cantidad            float64             `json:"cantidad"`
	Producto            Producto            `json:"producto"`
	Descripcion         string              `json:"descripcion"`
	PrecioUnitario      float64             `json:"precio_unitario"`
	Descuento           float64             `json:"descuento"`
	PorcentajeIVA       float64             `json:"porcentaje_IVA"`
	ValorIVA            float64             `json:"valor_IVA"`
	Total               float64             `json:"total"`
	EstadoOrdenServicio EstadoOrdenServicio `json:"estado_orden_servicio"`
}

func CreateResponseDetalleOrdenServicio(detalleOrdenServicioModel models.DetalleOrdenServicio, ordenServicio OrdenServicio, producto Producto, estadoOrdenServicio EstadoOrdenServicio) DetalleOrdenServicio {
	return DetalleOrdenServicio{ID: detalleOrdenServicioModel.ID, OrdenServicio: ordenServicio, Cantidad: detalleOrdenServicioModel.Cantidad, Producto: producto, Descripcion: detalleOrdenServicioModel.Descripcion, PrecioUnitario: detalleOrdenServicioModel.PrecioUnitario, Descuento: detalleOrdenServicioModel.Descuento, PorcentajeIVA: detalleOrdenServicioModel.PorcentajeIVA, ValorIVA: detalleOrdenServicioModel.ValorIVA, Total: detalleOrdenServicioModel.Total, EstadoOrdenServicio: estadoOrdenServicio}
}

func CreateDetalleOrdenesServicio(c *fiber.Ctx) error {
	var detalleOrdenServicio models.DetalleOrdenServicio

	if err := c.BodyParser(&detalleOrdenServicio); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var ordenServicio models.OrdenServicio
	if err := FindOrdenServicio(detalleOrdenServicio.OrdenServicioRefer, &ordenServicio); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var cliente models.Cliente
	if err := FindCliente(ordenServicio.ClienteRefer, &cliente); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var tecnico models.Tecnico
	if err := FindTecnico(ordenServicio.TecnicoRefer, &tecnico); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var empresa models.Empresa
	if err := findEmpresa(ordenServicio.EmpresaRefer, &empresa); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var ciudad models.Ciudad
	if err := FindCiudad(empresa.CiudadRefer, &ciudad); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var provincia models.Provincia
	if err := findProvincia(ciudad.ProvinciaRefer, &provincia); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var producto models.Producto
	if err := FindProducto(detalleOrdenServicio.ProductoRefer, &producto); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	var estadoOrdenServicio models.EstadoOrdenServicio
	if err := FindEstadoOrdenServicio(detalleOrdenServicio.EstadoOrdenServicioRefer, &estadoOrdenServicio); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.DB.Create(&detalleOrdenServicio)

	responseProvincia := CreateResponseProvincia(provincia)
	responseCiudad := CreateResponseCiudad(ciudad, responseProvincia)
	responseEmpresa := CreateResponseEmpresa(empresa, responseCiudad)
	responseTecnico := CreateResponseTecnico(tecnico, responseCiudad, responseEmpresa)
	responseCliente := CreateResponseCliente(cliente)
	responseEstadoOrdenServicio := CreateResponseEstadoOrdenServicio(estadoOrdenServicio, responseEmpresa)
	responseOrdenServicio := CreateResponseOrdenServicio(ordenServicio, responseEmpresa, responseEstadoOrdenServicio, responseTecnico, responseCliente)
	responseProducto := CreateResponseProducto(producto, responseEmpresa)
	responseDetalleOrdenServicio := CreateResponseDetalleOrdenServicio(detalleOrdenServicio, responseOrdenServicio, responseProducto, responseEstadoOrdenServicio)

	return c.Status(200).JSON(responseDetalleOrdenServicio)
}

func GetDetalleOrdenesServicio(c *fiber.Ctx) error {
	detalleOrdenesServicio := []models.DetalleOrdenServicio{}
	database.DB.Find(&detalleOrdenesServicio)
	responseDetalleOrdenesServicio := []DetalleOrdenServicio{}

	for _, detalleOrdenServicio := range detalleOrdenesServicio {

		var ordenServicio models.OrdenServicio
		database.DB.Find(&ordenServicio, "id = ?", detalleOrdenServicio.OrdenServicioRefer)

		var empresa models.Empresa
		database.DB.Find(&empresa, "id = ?", ordenServicio.EmpresaRefer)

		var ciudad models.Ciudad
		database.DB.Find(&ciudad, "id = ?", empresa.CiudadRefer)

		var provincia models.Provincia
		database.DB.Find(&provincia, "id = ?", ciudad.ProvinciaRefer)

		var tecnico models.Tecnico
		database.DB.Find(&tecnico, "id = ?", ordenServicio.TecnicoRefer)

		var cliente models.Cliente
		database.DB.Find(&tecnico, "id = ?", ordenServicio.TecnicoRefer)

		var producto models.Producto
		database.DB.Find(&producto, "id = ?", detalleOrdenServicio.ProductoRefer)

		var estadoOrdenServicio models.EstadoOrdenServicio
		database.DB.Find(&estadoOrdenServicio, "id = ?", ordenServicio.EstadoOrdenServicioRefer)

		responseCiudad := CreateResponseCiudad(ciudad, CreateResponseProvincia(provincia))
		responseEmpresa := CreateResponseEmpresa(empresa, responseCiudad)
		responseTecnico := CreateResponseTecnico(tecnico, responseCiudad, responseEmpresa)
		responseCliente := CreateResponseCliente(cliente)
		responseEstadoOrdenServicio := CreateResponseEstadoOrdenServicio(estadoOrdenServicio, responseEmpresa)
		responseOrdenServicio := CreateResponseOrdenServicio(ordenServicio, responseEmpresa, responseEstadoOrdenServicio, responseTecnico, responseCliente)
		responseProducto := CreateResponseProducto(producto, responseEmpresa)
		responseDetalleOrdenServicio := CreateResponseDetalleOrdenServicio(detalleOrdenServicio, responseOrdenServicio, responseProducto, responseEstadoOrdenServicio)
		responseDetalleOrdenesServicio = append(responseDetalleOrdenesServicio, responseDetalleOrdenServicio)
	}
	return c.Status(200).JSON(responseDetalleOrdenesServicio)
}

func FindDetalleOrdenServicio(id int, detalleOrdenServicio *models.DetalleOrdenServicio) error {
	database.DB.Find(&detalleOrdenServicio, "id = ?", id)
	if detalleOrdenServicio.ID == 0 {
		return errors.New("Order does not exist")
	}
	return nil
}

func GetDetalleOrdenServicio(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	var detalleOrdenServicio models.DetalleOrdenServicio

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := FindDetalleOrdenServicio(id, &detalleOrdenServicio); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	var ordenServicio models.OrdenServicio
	database.DB.First(&ordenServicio, detalleOrdenServicio.OrdenServicio)

	var empresa models.Empresa
	database.DB.First(&empresa, ordenServicio.EmpresaRefer)

	var ciudad models.Ciudad
	database.DB.First(&ciudad, empresa.CiudadRefer)

	var provincia models.Provincia
	database.DB.First(&provincia, ciudad.ProvinciaRefer)

	var tecnico models.Tecnico
	database.DB.First(&tecnico, ordenServicio.TecnicoRefer)

	var cliente models.Cliente
	database.DB.First(&cliente, ordenServicio.ClienteRefer)

	var producto models.Producto
	database.DB.First(&producto, detalleOrdenServicio.ProductoRefer)

	var estadoOrdenServicio models.EstadoOrdenServicio
	database.DB.First(&estadoOrdenServicio, detalleOrdenServicio.EstadoOrdenServicioRefer)

	responseProvincia := CreateResponseProvincia(provincia)
	responseCiudad := CreateResponseCiudad(ciudad, responseProvincia)
	responseEmpresa := CreateResponseEmpresa(empresa, responseCiudad)
	responseTecnico := CreateResponseTecnico(tecnico, responseCiudad, responseEmpresa)
	responseCliente := CreateResponseCliente(cliente)
	responseEstadoOrdenServicio := CreateResponseEstadoOrdenServicio(estadoOrdenServicio, responseEmpresa)
	responseProducto := CreateResponseProducto(producto, responseEmpresa)
	responseOrdenServicio := CreateResponseOrdenServicio(ordenServicio, responseEmpresa, responseEstadoOrdenServicio, responseTecnico, responseCliente)
	responseDetalleOrdenServicio := CreateResponseDetalleOrdenServicio(detalleOrdenServicio, responseOrdenServicio, responseProducto, responseEstadoOrdenServicio)

	return c.Status(200).JSON(responseDetalleOrdenServicio)
}

func FindDetalleOrdenServicioByOrder(orden_servicio_id int, detalleOrdenServicio *models.DetalleOrdenServicio) error {
	database.DB.Find(&detalleOrdenServicio, "orden_servicio_refer = ?", orden_servicio_id)
	if detalleOrdenServicio.OrdenServicioRefer == 0 {
		return errors.New("Order does not exist")
	}
	return nil
}

func GetDetalleOrdenServicioByOrder(c *fiber.Ctx) error {
	orden_servicio_id, err := c.ParamsInt("orden_servicio_id")

	detalleOrdenesServicio := []models.DetalleOrdenServicio{}
	database.DB.Find(&detalleOrdenesServicio)
	responseDetalleOrdenesServicio := []DetalleOrdenServicio{}

	for _, detalleOrdenServicio := range detalleOrdenesServicio {

		if err != nil {
			return c.Status(400).JSON("Please ensure that :idOrden is an integer")
		}

		if err := FindDetalleOrdenServicioByOrder(orden_servicio_id, &detalleOrdenServicio); err != nil {
			return c.Status(400).JSON(err.Error())
		}

		var ordenServicio models.OrdenServicio
		database.DB.Find(&ordenServicio, "id = ?", detalleOrdenServicio.OrdenServicioRefer)

		var empresa models.Empresa
		database.DB.Find(&empresa, "id = ?", ordenServicio.EmpresaRefer)

		var ciudad models.Ciudad
		database.DB.Find(&ciudad, "id = ?", empresa.CiudadRefer)

		var provincia models.Provincia
		database.DB.Find(&provincia, "id = ?", ciudad.ProvinciaRefer)

		var tecnico models.Tecnico
		database.DB.Find(&tecnico, "id = ?", ordenServicio.TecnicoRefer)

		var cliente models.Cliente
		database.DB.Find(&cliente, "id = ?", ordenServicio.ClienteRefer)

		var producto models.Producto
		database.DB.Find(&producto, "id = ?", detalleOrdenServicio.ProductoRefer)

		var estadoOrdenServicio models.EstadoOrdenServicio
		database.DB.Find(&estadoOrdenServicio, "id = ?", ordenServicio.EstadoOrdenServicioRefer)

		responseCiudad := CreateResponseCiudad(ciudad, CreateResponseProvincia(provincia))
		responseEmpresa := CreateResponseEmpresa(empresa, responseCiudad)
		responseTecnico := CreateResponseTecnico(tecnico, responseCiudad, responseEmpresa)
		responseCliente := CreateResponseCliente(cliente)
		responseEstadoOrdenServicio := CreateResponseEstadoOrdenServicio(estadoOrdenServicio, responseEmpresa)
		responseOrdenServicio := CreateResponseOrdenServicio(ordenServicio, responseEmpresa, responseEstadoOrdenServicio, responseTecnico, responseCliente)
		responseProducto := CreateResponseProducto(producto, responseEmpresa)
		responseDetalleOrdenServicio := CreateResponseDetalleOrdenServicio(detalleOrdenServicio, responseOrdenServicio, responseProducto, responseEstadoOrdenServicio)
		responseDetalleOrdenesServicio = append(responseDetalleOrdenesServicio, responseDetalleOrdenServicio)
	}
	return c.Status(200).JSON(responseDetalleOrdenesServicio)
}

func UpdatetDetalle(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var detalle models.DetalleOrdenServicio

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := FindDetalleOrdenServicio(id, &detalle); err != nil {
		return c.Status(400).JSON(err.Error())
	}
	type UpdateDetalleOrdenServicio struct {
		OrdenServicioRefer       int     `json:"orden_servicio_id"`
		Cantidad                 float64 `json:"cantidad"`
		ProductoRefer            int     `json:"producto_id"`
		Descripcion              string  `json:"descripcion"`
		PrecioUnitario           float64 `json:"precio_unitario"`
		Descuento                float64 `json:"descuento"`
		PorcentajeIVA            float64 `json:"porcentaje_IVA"`
		ValorIVA                 float64 `json:"valor_IVA"`
		Total                    float64 `json:"total"`
		EstadoOrdenServicioRefer int     `json:"estado_orden_servicio_id"`
	}

	var updateData UpdateDetalleOrdenServicio

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	detalle.OrdenServicioRefer = updateData.OrdenServicioRefer
	detalle.Cantidad = updateData.Cantidad
	detalle.ProductoRefer = updateData.ProductoRefer
	detalle.Descripcion = updateData.Descripcion
	detalle.PrecioUnitario = updateData.PrecioUnitario
	detalle.Descuento = updateData.Descuento
	detalle.PorcentajeIVA = updateData.PorcentajeIVA
	detalle.ValorIVA = updateData.ValorIVA
	detalle.Total = updateData.Total
	detalle.EstadoOrdenServicioRefer = updateData.EstadoOrdenServicioRefer

	var ordenServicio models.OrdenServicio
	database.DB.First(&ordenServicio, detalle.OrdenServicio)

	var empresa models.Empresa
	database.DB.First(&empresa, ordenServicio.EmpresaRefer)

	var ciudad models.Ciudad
	database.DB.First(&ciudad, empresa.CiudadRefer)

	var provincia models.Provincia
	database.DB.First(&provincia, ciudad.ProvinciaRefer)

	var tecnico models.Tecnico
	database.DB.First(&tecnico, ordenServicio.TecnicoRefer)

	var cliente models.Cliente
	database.DB.First(&cliente, ordenServicio.ClienteRefer)

	var producto models.Producto
	database.DB.First(&producto, detalle.ProductoRefer)

	var estadoOrdenServicio models.EstadoOrdenServicio
	database.DB.First(&estadoOrdenServicio, detalle.EstadoOrdenServicioRefer)

	database.DB.Save(&detalle)

	responseProvincia := CreateResponseProvincia(provincia)
	responseCiudad := CreateResponseCiudad(ciudad, responseProvincia)
	responseEmpresa := CreateResponseEmpresa(empresa, responseCiudad)
	responseTecnico := CreateResponseTecnico(tecnico, responseCiudad, responseEmpresa)
	responseCliente := CreateResponseCliente(cliente)
	responseEstadoOrdenServicio := CreateResponseEstadoOrdenServicio(estadoOrdenServicio, responseEmpresa)
	responseProducto := CreateResponseProducto(producto, responseEmpresa)
	responseOrdenServicio := CreateResponseOrdenServicio(ordenServicio, responseEmpresa, responseEstadoOrdenServicio, responseTecnico, responseCliente)
	responseDetalleOrdenServicio := CreateResponseDetalleOrdenServicio(detalle, responseOrdenServicio, responseProducto, responseEstadoOrdenServicio)

	return c.Status(200).JSON(responseDetalleOrdenServicio)
}
