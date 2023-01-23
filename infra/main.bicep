param location string
param acrName string
param tag string

param resizeApiName string
param resizeApiPort string

// param storeApiName string
// param storeApiPort string

// param photoApiName string
// param photoApiPort string

param uploadsStorageQueueName string
param imagesStorageQueueName string
param thumbsContainerName string
param imagesContainerName string
param uploadsContainerName string

param databaseName string
param containerName string
param partitionKey string

param maxThumbHeight string
param maxThumbWidth string
param maxImageHeight string
param maxImageWidth string

param grpcMaxRequestSizeMb int

param tags object = {
  environment: 'dev'
  costcode: '1234567890'
}

var affix = uniqueString(resourceGroup().id)

var containerAppEnvName = 'app-env-${affix}'
var resizeApiContainerImage = '${acr.properties.loginServer}/${resizeApiName}:${tag}'
//var storeApiContainerImage = '${acr.properties.loginServer}/${storeApiName}:${tag}'
//var photoApiContainerImage = '${acr.properties.loginServer}/${photoApiName}:${tag}'
var acrLoginServer = '${acrName}.azurecr.io'
var acrAdminPassword = listCredentials(acr.id, '2021-12-01-preview').passwords[0].value
var workspaceName = 'wks-${affix}'
var storageAccountName = 'stor${affix}'
var aiName = 'ai-${affix}'
var secrets = [
  {
    name: 'registry-password'
    value: acrAdminPassword
  }
]

module aiModule 'modules/ai.bicep' = {
  name: 'module-ai'
  params: {
    location: location
    aiName: aiName
  }
}

resource acr 'Microsoft.ContainerRegistry/registries@2021-12-01-preview' existing = {
  name: acrName
}

module wksModule 'modules/wks.bicep' = {
  name: 'module-wks'
  params: {
    location: location
    name: workspaceName
    tags: tags
  }
}

module storModule 'modules/stor.bicep' = {
  name: 'module-stor'
  params: {
    location: location
    name: storageAccountName
    containers: [
      {
        name: thumbsContainerName
        publicAccess: 'Blob'
      }
      {
        name: imagesContainerName
        publicAccess: 'Blob'
      }
      {
        name: uploadsContainerName
        publicAccess: 'None'
      }
/*       thumbsContainerName
      imagesContainerName
      uploadsContainerName */
    ]
    tags: tags
  }
}

module cosmosModule 'modules/cosmos.bicep' = {
  name: 'module-cosmos'
  params: {
    location: location
    tags: tags
    containerName: containerName
    partitionKey: partitionKey
    databaseName: databaseName
  }
}

module eventGridModule 'modules/eventgrid.bicep' = {
  name: 'module-evg'
  params: {
    name: 'egt'
    containers: [
      uploadsContainerName
      imagesContainerName
      thumbsContainerName
    ]
    storageAccountId: storModule.outputs.id
    location: location
    eventTypes: [
      'Microsoft.Storage.BlobCreated'
    ]
    topicSourceId: storModule.outputs.id
    tags: tags
  }
}

module containerAppEnvModule './modules/cappenv.bicep' = {
  name: 'module-capp-env'
  params: {
    name: containerAppEnvName
    location: location
    isInternal: false
    tags: tags
    wksSharedKey: wksModule.outputs.workspaceSharedKey
    aiKey: aiModule.outputs.aiKey
    wksCustomerId: wksModule.outputs.workspaceCustomerId
  }
}

resource resizeApi 'Microsoft.App/containerApps@2022-06-01-preview' = {
  name: resizeApiName
  location: location
  tags: tags
  identity: {
    type: 'SystemAssigned'
  }
  dependsOn: [
    containerAppEnvModule
    storModule
    eventGridModule
  ]
  properties: {
    configuration: {
      activeRevisionsMode: 'single'
      dapr: {
        appId: resizeApiName
        appPort: int(resizeApiPort)
        appProtocol: 'grpc'
        enabled: true
        httpMaxRequestSize: grpcMaxRequestSizeMb
      }
      secrets: secrets
      registries: [
        {
          passwordSecretRef: 'registry-password'
          server: acrLoginServer
          username: acr.name
        }
      ]
      ingress: {
        external: false
        targetPort: int(resizeApiPort)
        traffic: [
          {
            latestRevision: true
            weight: 100
          }
        ]
        transport: 'http'
      }
    }
    managedEnvironmentId: containerAppEnvModule.outputs.id
    template: {
      containers: [
        {
          image: resizeApiContainerImage
          name: resizeApiName
          resources: {
            cpu: '0.25'
            memory: '0.5Gi'
          }
          env: [
            {
              name: 'SERVICE_NAME'
              value: resizeApiName
            }
            {
              name: 'SERVICE_PORT'
              value: resizeApiPort
            }
            {
              name: 'UPLOADS_QUEUE_BINDING'
              value: 'queue-${uploadsContainerName}'
            }
            {
              name: 'IMAGES_CONTAINER_BINDING'
              value: 'blob-${imagesContainerName}'
            }
            {
              name: 'UPLOADS_CONTAINER_BINDING'
              value: 'blob-${uploadsContainerName}'
            }
            {
              name: 'THUMBS_CONTAINER_BINDING'
              value: 'blob-${thumbsContainerName}'
            }
            {
              name: 'MAX_THUMB_HEIGHT'
              value: maxThumbHeight
            }
            {
              name: 'MAX_THUMB_WIDTH'
              value: maxThumbWidth
            }
            {
              name: 'MAX_IMAGE_HEIGHT'
              value: maxImageHeight
            }
            {
              name: 'MAX_IMAGE_WIDTH'
              value: maxImageWidth
            }
            {
              name: 'GRPC_MAX_REQUEST_BODY_SIZE_MB'
              value: string(grpcMaxRequestSizeMb)
            }
            {
              name: 'STORAGE_ACCOUNT_NAME'
              value: storModule.outputs.name
            }
          ]
        }
      ]
      scale: {
        minReplicas: 1
        maxReplicas: 10
      }
    }
  }
}

/* resource storeApi 'Microsoft.App/containerApps@2022-10-01' = {
  name: storeApiName
  location: location
  tags: tags
  identity: {
    type: 'SystemAssigned'
  }
  dependsOn: [
    containerAppEnvModule
    cosmosModule
  ]
  properties: {
    configuration: {
      activeRevisionsMode: 'single'
      dapr: {
        appId: storeApiName
        appPort: int(storeApiPort)
        appProtocol: 'grpc'
        enabled: true
        httpMaxRequestSize: grpcMaxRequestSizeMb
      }
      secrets: secrets
      registries: [
        {
          passwordSecretRef: 'registry-password'
          server: acrLoginServer
          username: acr.name
        }
      ]
      ingress: {
        external: false
        targetPort: int(storeApiPort)
        traffic: [
          {
            latestRevision: true
            weight: 100
          }
        ]
        transport: 'http'
      }
    }
    managedEnvironmentId: containerAppEnvModule.outputs.id
    template: {
      containers: [
        {
          image: storeApiContainerImage
          name: storeApiName
          resources: {
            cpu: '0.25'
            memory: '0.5Gi'
          }
          env: [
            {
              name: 'SERVICE_NAME'
              value: storeApiName
            }
            {
              name: 'SERVICE_PORT'
              value: storeApiPort
            }
            {
              name: 'COSMOS_BINDING'
              value: 'db-${databaseName}'
            }
            {
              name: 'IMAGES_QUEUE_BINDING'
              value: 'queue-${imagesContainerName}'
            }
            {
              name: 'IMAGES_CONTAINER_BINDING'
              value: 'blob-${imagesContainerName}'
            }
            {
              name: 'UPLOADS_CONTAINER_BINDING'
              value: 'blob-${uploadsContainerName}'
            }
            {
              name: 'THUMBS_CONTAINER_BINDING'
              value: 'blob-${thumbsContainerName}'
            }
          ]
        }
      ]
      scale: {
        minReplicas: 1
        maxReplicas: 10
      }
    }
  }
}

resource photoApi 'Microsoft.App/containerApps@2022-10-01' = {
  name: photoApiName
  location: location
  tags: tags
  identity: {
    type: 'SystemAssigned'
  }
  dependsOn: [
    containerAppEnvModule
    cosmosModule
  ]
  properties: {
    configuration: {
      activeRevisionsMode: 'single'
      dapr: {
        appId: photoApiName
        appPort: int(photoApiPort)
        appProtocol: 'http'
        enabled: true
        httpMaxRequestSize: grpcMaxRequestSizeMb
      }
      secrets: secrets
      registries: [
        {
          passwordSecretRef: 'registry-password'
          server: acrLoginServer
          username: acr.name
        }
      ]
      ingress: {
        external: true
        targetPort: int(photoApiPort)
        traffic: [
          {
            latestRevision: true
            weight: 100
          }
        ]
        transport: 'http'
        corsPolicy: {
          allowedOrigins: [
            '*'
          ]
          allowedMethods: [
            'GET'
            'POST'
            'OPTIONS'
            'DELETE'
            'PUT'
          ]
          allowedHeaders: [
            '*'
          ]
          allowCredentials: true
       }
      }
    }
    managedEnvironmentId: containerAppEnvModule.outputs.id
    template: {
      containers: [
        {
          image: photoApiContainerImage
          name: photoApiName
          resources: {
            cpu: '0.25'
            memory: '0.5Gi'
          }
          env: [
            {
              name: 'SERVICE_NAME'
              value: photoApiName
            }
            {
              name: 'SERVICE_PORT'
              value: photoApiPort
            }
            {
              name: 'COSMOS_BINDING'
              value: 'db-${databaseName}'
            }
            {
              name: 'COSMOS_URL'
              value: cosmosModule.outputs.url
            }
            {
              name: 'COSMOS_DB'
              value: cosmosModule.outputs.dbName
            }
            {
              name: 'COSMOS_KEY'
              value: cosmosModule.outputs.key
            }
            {
              name: 'COSMOS_CONTAINER'
              value: cosmosModule.outputs.containerName
            }
          ]
        }
      ]
      scale: {
        minReplicas: 1
        maxReplicas: 10
      }
    }
  }
} */

resource uploadsStorageQueueDaprComponent 'Microsoft.App/managedEnvironments/daprComponents@2022-06-01-preview' = {
  dependsOn: [
    containerAppEnvModule
  ]
  name: '${containerAppEnvName}/queue-${toLower(uploadsStorageQueueName)}'
  properties: {
    componentType: 'bindings.azure.storagequeues'
    version: 'v1'
    ignoreErrors: false
    initTimeout: '60s'
    metadata: [
      {
        name: 'storageAccount'
        value: storModule.outputs.name
      }
      {
        name: 'storageAccessKey'
        value: storModule.outputs.key
      }
      {
        name: 'queue'
        value: uploadsStorageQueueName
      }
    ]
    scopes: [
      resizeApiName
    ]
  }
}

resource imagesStorageQueueDaprComponent 'Microsoft.App/managedEnvironments/daprComponents@2022-06-01-preview' = {
  dependsOn: [
    containerAppEnvModule
  ]
  name: '${containerAppEnvName}/queue-${toLower(imagesStorageQueueName)}'
  properties: {
    componentType: 'bindings.azure.storagequeues'
    version: 'v1'
    ignoreErrors: false
    initTimeout: '60s'
    metadata: [
      {
        name: 'storageAccount'
        value: storModule.outputs.name
      }
      {
        name: 'storageAccessKey'
        value: storModule.outputs.key
      }
      {
        name: 'queue'
        value: imagesStorageQueueName
      }
    ]
    scopes: [
      resizeApiName
    ]
  }
}

resource uploadsStorageDaprComponent 'Microsoft.App/managedEnvironments/daprComponents@2022-06-01-preview' = {
  dependsOn: [
    containerAppEnvModule
  ]
  name: '${containerAppEnvName}/blob-${toLower(uploadsContainerName)}'
  properties: {
    componentType: 'bindings.azure.blobstorage'
    version: 'v1'
    ignoreErrors: false
    initTimeout: '60s'
    metadata: [
      {
        name: 'storageAccount'
        value: storModule.outputs.name
      }
      {
        name: 'storageAccessKey'
        value: storModule.outputs.key
      }
      {
        name: 'container'
        value: uploadsContainerName
      }
    ]
    scopes: [
      resizeApiName
    ]
  }
}

resource thumbsStorageDaprComponent 'Microsoft.App/managedEnvironments/daprComponents@2022-06-01-preview' = {
  dependsOn: [
    containerAppEnvModule
  ]
  name: '${containerAppEnvName}/blob-${toLower(thumbsContainerName)}'
  properties: {
    componentType: 'bindings.azure.blobstorage'
    version: 'v1'
    ignoreErrors: false
    initTimeout: '60s'
    metadata: [
      {
        name: 'storageAccount'
        value: storModule.outputs.name
      }
      {
        name: 'storageAccessKey'
        value: storModule.outputs.key
      }
      {
        name: 'container'
        value: thumbsContainerName
      }
    ]
    scopes: [
      resizeApiName
    ]
  }
}

resource imagesStorageDaprComponent 'Microsoft.App/managedEnvironments/daprComponents@2022-06-01-preview' = {
  dependsOn: [
    containerAppEnvModule
  ]
  name: '${containerAppEnvName}/blob-${toLower(imagesContainerName)}'
  properties: {
    componentType: 'bindings.azure.blobstorage'
    version: 'v1'
    ignoreErrors: false
    initTimeout: '60s'
    metadata: [
      {
        name: 'storageAccount'
        value: storModule.outputs.name
      }
      {
        name: 'storageAccessKey'
        value: storModule.outputs.key
      }
      {
        name: 'container'
        value: imagesContainerName
      }
    ]
    scopes: [
      resizeApiName
    ]
  }
}

/* resource cosmosDbDaprComponent 'Microsoft.App/managedEnvironments/daprComponents@2022-06-01-preview' = {
  dependsOn: [
    containerAppEnvModule
  ]
  name: '${containerAppEnvName}/db-${toLower(databaseName)}'
  properties: {
    componentType: 'state.azure.cosmosdb'
    version: 'v1'
    ignoreErrors: false
    initTimeout: '60s'
    metadata: [
      {
        name: 'url'
        value: cosmosModule.outputs.url
      }
      {
        name: 'masterKey'
        value: cosmosModule.outputs.key
      }
      {
        name: 'database'
        value: cosmosModule.outputs.dbName
      }
      {
        name: 'collection'
        value: cosmosModule.outputs.containerName
      }
    ]
    scopes: [
      storeApiName
      photoApiName
    ]
  }
} */

output resizeUrl string = resizeApi.properties.configuration.ingress.fqdn
output storageAccount string = storModule.outputs.name
