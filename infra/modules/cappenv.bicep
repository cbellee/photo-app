param location string
param name string
param aiKey string
param wksCustomerId string
param wksSharedKey string
param tags object
param acaInfraSubnetId string = ''
param acaRuntimSubnetId string = ''
param vnetConfig object = {
  internal: isInternal
  infrastructureSubnetId: acaInfraSubnetId
  runtimeSubnetId: acaRuntimSubnetId
  platformReservedCidr: '10.0.0.0/16'
  platformReservedDnsIP: '10.0.0.2'
  dockerBridgeCidr: '10.1.0.1/16'
}
param isInternal bool = true
param isZoneRedundant bool = false

resource containerAppEnvironment 'Microsoft.App/managedEnvironments@2022-03-01' = {
  location: location
  name: name
  properties: {
    daprAIInstrumentationKey: aiKey
    vnetConfiguration: isInternal ? vnetConfig : null
    zoneRedundant: isZoneRedundant
    appLogsConfiguration: {
      destination: 'log-analytics'
      logAnalyticsConfiguration: {
        customerId: wksCustomerId
        sharedKey: wksSharedKey
      }
    }
  }
  tags: tags
}

output id string = containerAppEnvironment.id
output name string = containerAppEnvironment.name
output cAppEnv object = containerAppEnvironment
