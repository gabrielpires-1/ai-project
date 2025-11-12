variable "resource_group_name" {
  description = "Nome do Grupo de Recursos"
  type        = string
  default     = "rg-gemini-api"
}

variable "location" {
  description = "Região onde os recursos serão criados"
  type        = string
  default     = "East US"
}

variable "acr_name_prefix" {
  description = "Prefixo para o nome do Container Registry (será adicionado um sufixo único)"
  type        = string
  default     = "acrgeminiapi"
}

variable "aca_app_name" {
  description = "Nome da sua aplicação no Container Apps"
  type        = string
  default     = "api-gemini-app"
}