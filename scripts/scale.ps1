param (
    $storageAccountName = 'storaodiwusgqeyiu',
    $containerName = 'uploads',
    $localPath = '../images',
    $rgName = 'photo-app-test-rg',
    $containers = @('uploads','thumbs','images')
)

# Install-Module -Name Az.App
# Register-AzResourceProvider -ProviderNamespace Microsoft.App

# remove existing blobs
$ctx = (Get-AzStorageAccount -ResourceGroupName $rgName -Name $storageAccountName).Context 
$containers | ForEach-Object {
    Get-AzStorageBlob -Container $_ -Context $ctx | Remove-AzStorageBlob 
}

# list # of replicas
"# of Replicas: $(Get-AzContainerAppRevision -ContainerAppName resize -ResourceGroupName $rgName | 
    Where-Object Active -eq $true | 
        Select-Object -ExpandProperty replica)"

Get-ChildItem $localPath -Directory | ForEach-Object {
    $collection = $_.Name
    Get-ChildItem $_.FullName -Directory | ForEach-Object {
        $album = $_.Name
        Get-ChildItem $_.FullName -File | ForEach-Object {
        Set-AzStorageBlobContent `
            -File $($_.FullName) `
            -Container $containerName `
            -Blob "$collection/$album/$($_.Name)" `
            -BlobType Block `
            -Properties @{'ContentType'='image/jpeg'} `
            -Metadata @{collection=$collection;album=$album;name=$($_.Name)} `
            -Tag @{collection=$collection;album=$album} `
            -Context $ctx `
            -Force
        }
<#         azcopy copy "$($_.FullName)/*" "https://$storageAccountName.blob.core.windows.net/$containerName/$collection/$album" `
        --blob-tags="collection=$collection&album=$album" `
        --metadata "collection=$collection;album=$album;name=$($_.Name)" #>
    }
}

# list # of replicas
"# of Replicas: $(Get-AzContainerAppRevision -ContainerAppName resize -ResourceGroupName $rgName | 
    Where-Object Active -eq $true | 
        Select-Object -ExpandProperty replica)"
