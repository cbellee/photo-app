import { BlockBlobClient } from "@azure/storage-blob";
import azureBlobStorageSender from "../services/azureBlobStorageSender";
import { apiConfig, msalConfig } from "../authConfig";
import { InteractiveBrowserCredential } from "@azure/identity"

const [collectionAlbumData, setCollectionAlbumData] = useState([]);
const [albumData, setAlbumData] = useState([]);
const [currentAlbum, setCurrentAlbum] = useState([]);
const [currentCollection, setCurrentCollection] = useState([]);

var signInOptions = {
  clientId: msalConfig.auth.clientId,
  tenantId: apiConfig.tenantId
}

var browserCredential = new InteractiveBrowserCredential(signInOptions);

const UploadButtonWithEvents = () => {
  useItemFinishListener((file) => {
    setPhotoName(file.uploadResponse.path);
    setOriginalPhoto(file.uploadResponse.publicUrl);
  });

  return <UploadButton />;
};

let client = new BlockBlobClient(url, browserCredential)

export const UploadBlob = () => (
  <>
  <Uploady
    multiple={false}
    send={azureBlobStorageSender({
      blobClient: blobClient,
      container: "rooms",
    })}
  >
    <UploadButtonWithEvents />
    <UploadPreview />
  </Uploady>


  <form onSubmit={(ev) => ev.target.reset()}>
    <div className="inputs">
      <p>
        <input
          type="file"
          id="file"
          onChange={changeHandler}
          accept="image/png, image/jpg, image/jpeg"
          multiple
        />
      </p>

      {
        imageFiles.length > 0 ?

          <input
            type="button"
            id="upload"
            name="upload"
            value="upload"
            className="inputs"
            disabled={isUploading}
            onClick={sendRequest}
          />
          : null
      }
      <InputLabel>Collections</InputLabel>
      <CreatableSelect
        onChange={(event, value) => handleCollection(event, value)}
        isClearable
        options={
          Array.from(collectionAlbumData.keys()).map((option, idx) => {
            return { value: option, label: option }
          })
        }
      >
      </CreatableSelect>
      <InputLabel>Albums</InputLabel>
      <CreatableSelect
        name="Albums"
        onChange={(event, value) => handleAlbum(event, value)}
        isClearable
        options={
          albumData.map((option, idx) => {
            return { value: idx, label: option }
          })
        }
      >
      </CreatableSelect>
    </div>
    <div className={classes.imageList}>
      {
        images.length > 0 ?
          <ImageList spacing={30} className={classes.imageList}>
            <ImageListItem key="Subheader" cols={4} style={{ height: 'auto' }}>
              <ListSubheader component="div"></ListSubheader>
            </ImageListItem>
            {images.map((image, idx) => {
              return <>
                <ImageListItem key={idx} className={classes.imageItem}>
                  <img src={image.blob} alt={image.blob} />
                  <ImageListItemBar
                    className={classes.imageTitleBar}
                    title={`${image.name} (width ${image.width} x height ${image.height})`}
                    actionIcon={<IconButton aria-label={`info about ${image.name}`} className={classes.icon}>
                      <FormControl>
                        <RadioGroup row name="album-thumbnail-group" value="">
                          <FormControlLabel
                            key={idx}
                            value={idx}
                            checked={isAlbumThumbnail[idx] === 'true'}
                            control={<Radio />}
                            onChange={updateRadioThumbnailChanged(idx)}
                            label="album thumbnail" />
                        </RadioGroup>
                      </FormControl>
                    </IconButton>} />
                </ImageListItem>
              </>
            })}
          </ImageList> : null
      }
    </div>
  </form>
  </>
);


useEffect(() => {
  getAlbumCollections("uploads")
    .catch(error => console.log(error), (data = '[]') => console.log("data: " + data))
    .then((data) => { setCollectionAlbumData(data) });
}, [])

const handleCollection = (event, value) => {
  let albums = collectionAlbumData.get(event.value);
  if (albums && albums.length > 0) {
    setAlbumData(albums);
    setCurrentCollection(event.value, [currentCollection]);
    console.log("current collection: " + event.value)
  } else {
    setAlbumData([]);
    setCurrentCollection(event.value, [currentAlbum]);
    console.log("current collection: " + event.value)
  }
}

const handleAlbum = (event, value) => {
  setCurrentAlbum(event.label);
  console.log("current album: " + event.label)
}
