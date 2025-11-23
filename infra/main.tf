# Adiciona o módulo 'random' para garantir nomes únicos (necessário para ACR)
resource "random_id" "suffix" {
  byte_length = 6
}

# 1. Grupo de Recursos (Pré-requisito para todos os recursos)
resource "azurerm_resource_group" "main" {
  name     = var.resource_group_name
  location = var.location
}

# 2. Azure Container Registry (ACR)
resource "azurerm_container_registry" "main" {
  # Combina o prefixo da variável com o sufixo aleatório para garantir nome único
  name                = "${var.acr_name_prefix}${random_id.suffix.hex}"
  resource_group_name = azurerm_resource_group.main.name
  location            = azurerm_resource_group.main.location
  
  # SKU "Basic" e "admin_enabled = true" para simplicidade e CI/CD
  sku                 = "Basic"
  admin_enabled       = true 
}

# 3. Ambiente do Container Apps (A "vizinhança" onde o App roda)
resource "azurerm_log_analytics_workspace" "main" {
  name                = "logs-${random_id.suffix.hex}"
  location            = azurerm_resource_group.main.location
  resource_group_name = azurerm_resource_group.main.name
  sku                 = "PerGB2018"
}

resource "azurerm_container_app_environment" "main" {
  name                       = "env-${random_id.suffix.hex}"
  location                   = azurerm_resource_group.main.location
  resource_group_name        = azurerm_resource_group.main.name
  log_analytics_workspace_id = azurerm_log_analytics_workspace.main.id
}

# 4. Azure Container App (ACA) - (Sua API rodando)
resource "azurerm_container_app" "main" {
  name                         = var.aca_app_name
  resource_group_name          = azurerm_resource_group.main.name
  container_app_environment_id = azurerm_container_app_environment.main.id

  revision_mode = "Single"

  # 1. Definição dos 'secrets'
  secret {
    name  = "acr-password"
    value = azurerm_container_registry.main.admin_password
  }
  secret {
    name = "google-api-key"
    value = "PLACEHOLDER_KEY_PARA_TERRAFORM" 
  }

  # 2. Configuração do registro de imagem.
  registry {
    server               = azurerm_container_registry.main.login_server
    username             = azurerm_container_registry.main.admin_username
    password_secret_name = "acr-password" 
  }

  template {
    container {
      name   = "api-gemini"
      image  = "mcr.microsoft.com/azuredocs/containerapps-helloworld:latest" 
      cpu    = 0.25
      memory = "0.5Gi"

      env {
        name        = "GOOGLE_API_KEY"
        secret_name = "google-api-key" 
      }

      env {
        name  = "PORT"
        value = "8080"
      }
    }
  }

  # Configuração do Ingress (Endpoint Público)
  ingress {
    external_enabled = true  
    # Balanceador de carga do Azure Container Apps vai enviar tráfego para a porta 8080 do container
    target_port      = 8080 
    
    traffic_weight {
      percentage      = 100
      latest_revision = true
    }
  }
}

# Saída do URL da API para o deploy e acesso
output "container_app_url" {
  description = "O URL público para acessar a API"
  value       = azurerm_container_app.main.ingress[0].fqdn
}

# Saídas úteis para o próximo passo (GitHub Actions)
output "acr_login_server" {
  description = "Servidor de login do ACR (ex: acrgeminiapi1234.azurecr.io)"
  value       = azurerm_container_registry.main.login_server
}

output "acr_full_name" {
  description = "Nome completo gerado do ACR"
  value       = azurerm_container_registry.main.name
}