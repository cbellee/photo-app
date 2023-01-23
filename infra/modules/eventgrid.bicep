param location string
param name string
param topicSourceId string
param tags object
param containers array
param storageAccountId string

@allowed([
  'Microsoft.Storage.BlobCreated'
  'Microsoft.Storage.BlobDeleted'
  'Microsoft.Storage.BlobTierChanged'
  'Microsoft.Storage.AsyncOperationInitiated'
])
param eventTypes array = [
  'Microsoft.Storage.BlobCreated'
]

var suffix = uniqueString(resourceGroup().id)
var topicName = '${name}-${suffix}'
var topicType = 'Microsoft.Storage.StorageAccounts'

resource eventGridTopic 'Microsoft.EventGrid/systemTopics@2021-12-01' = {
  location: location
  name: topicName
  identity: {
    type: 'SystemAssigned'
  }
  properties: {
    source: topicSourceId
    topicType: topicType
  }
  tags: tags
}

@batchSize(1)
resource eventGridSubscription 'Microsoft.EventGrid/systemTopics/eventSubscriptions@2021-12-01' = [for container in containers: {
  name: container
  parent: eventGridTopic
  properties: {
    destination: {
      endpointType: 'StorageQueue'
      properties: {
        queueName: container
        queueMessageTimeToLiveInSeconds: 600
        resourceId: storageAccountId
      }
    }
    eventDeliverySchema: 'EventGridSchema'
    retryPolicy: {
      eventTimeToLiveInMinutes: 5
      maxDeliveryAttempts: 10
    }
    filter: {
      subjectBeginsWith: '/blobServices/default/containers/${container}/'
      includedEventTypes: eventTypes
      enableAdvancedFilteringOnArrays: false
    }
  }
}]

output topicName string = eventGridTopic.name
output topicId string = eventGridTopic.id
