param databaseName string
param location string
param tags object
param containerName string
param partitionKey string = '/url'
param isFreeTier bool = false

var accountName = 'cosmos-${uniqueString(resourceGroup().id)}'

resource cosmosAccount 'Microsoft.DocumentDB/databaseAccounts@2021-10-15' = {
  name: toLower(accountName)
  location: location
  properties: {
    enableFreeTier: isFreeTier
    locations: [
      {
        locationName: location
      }
    ]
    databaseAccountOfferType: 'Standard'
  }
}

resource sqlDB 'Microsoft.DocumentDB/databaseAccounts/sqlDatabases@2021-10-15' = {
  name: toLower(databaseName)
  tags: tags
  parent: cosmosAccount
  properties: {
    resource: {
      id: databaseName
    }
    options: {
      throughput: 400
    }
  }
}

resource container 'Microsoft.DocumentDB/databaseAccounts/sqlDatabases/containers@2021-10-15' = {
  name: containerName
  parent: sqlDB
  tags: tags
  properties: {
    resource: {
      partitionKey: {
        paths: [
          partitionKey
        ]
      }
      id: containerName
    }
  }
}

output url string = cosmosAccount.properties.documentEndpoint
output key string = cosmosAccount.listKeys().primaryMasterKey
output accountName string = cosmosAccount.name
output dbName string = sqlDB.name
output containerName string = container.name
output partitionKey string = partitionKey
