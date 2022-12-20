param location string
param name string
param tags object
param containers array
param isPublicBlobAccessAllowed bool = true
param isSupportHttpsTrafficOnly bool = true
param isDefaultToOAuthAuthentication bool = false
param isAllowSharedAccessKey bool = true

@allowed([
  'Storage'
  'StorageV2'
  'BlobStorage'
  'BlockBlobStorage'
  'FileStorage'
])
param kind string = 'StorageV2'

@allowed([
  'Standard_LRS'
  'Premium_LRS'
  'Premium_ZRS'
])
param sku string = 'Standard_LRS'

@allowed([
  'Cool'
  'Hot'
  'Premium'
])
param accessTier string = 'Hot'

@allowed([
  'Enabled'
  'Disabled'
])
param isPublicNetworkAccessEnabled string = 'Enabled'

resource storage 'Microsoft.Storage/storageAccounts@2021-09-01' = {
  kind: kind
  location: location
  name: name
  sku: {
    name: sku
  }
  properties: {
    accessTier: accessTier
    allowBlobPublicAccess: isPublicBlobAccessAllowed
    defaultToOAuthAuthentication: isDefaultToOAuthAuthentication
    publicNetworkAccess: isPublicNetworkAccessEnabled
    supportsHttpsTrafficOnly: isSupportHttpsTrafficOnly
    allowSharedKeyAccess: isAllowSharedAccessKey
  }
  tags: tags
}

resource queueService 'Microsoft.Storage/storageAccounts/queueServices@2021-09-01' = {
  parent: storage
  name: 'default'
}

resource blobService 'Microsoft.Storage/storageAccounts/blobServices@2021-09-01' = {
  parent: storage
  name: 'default'
}

resource storageQueues 'Microsoft.Storage/storageAccounts/queueServices/queues@2021-09-01' = [for container in containers: {
  parent: queueService
  name: container.name
}]

resource blobContainers 'Microsoft.Storage/storageAccounts/blobServices/containers@2021-09-01' = [for container in containers: {
  parent: blobService
  name: container.name
  properties: {
    publicAccess: container.publicAccess
  }
}]

output name string = storage.name
output id string = storage.id
output key string = storage.listKeys().keys[0].value
