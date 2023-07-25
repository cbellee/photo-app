param location string = resourceGroup().location
param name string
param tags object
param retentionInDays int = 30

@allowed([
  'Standard'
  'Free'
  'Premium'
  'PerGB2018'
])
param sku string = 'PerGB2018'

resource wks 'Microsoft.OperationalInsights/workspaces@2021-06-01' = {
  location: location
  name: name
  tags: tags
  properties: {
    retentionInDays: retentionInDays
    sku: {
      name: sku
    }
  }
}

output workspaceId string = wks.id 
output workspaceSharedKey string = wks.listKeys().primarySharedKey
output workspaceCustomerId string = wks.properties.customerId
