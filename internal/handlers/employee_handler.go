package handlers

import (
	"laba6/internal/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetEmployees
// @Summary      Get all employees
// @Description  Returns list of all employees from database
// @Tags         employees
// @Accept       json
// @Produce      json
// @Success      200  {array}   models.Employee
// @Failure      500  {object}  map[string]string
// @Router       /v1/employees [get]
func (h *Handler) GetEmployees(c *gin.Context) {
	employees, err := h.processors.EmployeeProcessor.GetAllEmployees()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get employees"})
		return
	}
	c.JSON(http.StatusOK, employees)
}

// GetEmployee
// @Summary      Get employee by ID
// @Description  Returns single employee by ID
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Employee ID"
// @Success      200  {object}  models.Employee
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /v1/employees/{id} [get]
func (h *Handler) GetEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	employee, err := h.processors.EmployeeProcessor.GetEmployeeByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}
	c.JSON(http.StatusOK, employee)
}

// CreateEmployee
// @Summary      Create employee
// @Description  Adds a new employee
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        employee  body      models.Employee  true  "Employee"
// @Success      201  {object}  models.Employee
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /v1/employees [post]
func (h *Handler) CreateEmployee(c *gin.Context) {
	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	err := h.processors.EmployeeProcessor.CreateEmployee(
		employee.Name,
		employee.Position,
		employee.Department,
		employee.Salary,
	)
	if err != nil {
		// Check error type and return appropriate status
		errMsg := err.Error()

		if strings.Contains(errMsg, "already exists") {
			c.JSON(http.StatusConflict, gin.H{"error": errMsg})
			return
		}

		if strings.Contains(errMsg, "cannot be negative") ||
			strings.Contains(errMsg, "is required") {
			c.JSON(http.StatusBadRequest, gin.H{"error": errMsg})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create employee"})
		return
	}

	c.JSON(http.StatusCreated, employee)
}

// DeleteEmployee
// @Summary      Delete employee
// @Description  Deletes employee by ID
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Employee ID"
// @Success      204  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /v1/employees/{id} [delete]
func (h *Handler) DeleteEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	if err := h.processors.EmployeeProcessor.DeleteEmployee(id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	c.Status(http.StatusNoContent)
}

// UpdateEmployee
// @Summary      Update employee
// @Description  Updates employee by ID
// @Tags         employees
// @Accept       json
// @Produce      json
// @Param        id        path      int              true  "Employee ID"
// @Param        employee  body      models.Employee  true  "Employee"
// @Success      200  {object}  models.Employee
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /v1/employees/{id} [put]
func (h *Handler) UpdateEmployee(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid employee ID"})
		return
	}

	var employee models.Employee
	if err := c.ShouldBindJSON(&employee); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	err = h.processors.EmployeeProcessor.UpdateEmployee(id, employee.Name, employee.Position, employee.Department, employee.Salary)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Employee not found"})
		return
	}

	employee.ID = id
	c.JSON(http.StatusOK, employee)
}
