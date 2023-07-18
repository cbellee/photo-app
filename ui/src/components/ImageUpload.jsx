import { useEffect, useState, useRef, useCallback, React } from "react";
import ImageList from '@mui/material/ImageList';
import { makeStyles } from '@mui/styles';
import { useMsal } from "@azure/msal-react";
import ImageListItem from '@mui/material/ImageListItem';
import ImageListItemBar from '@mui/material/ImageListItemBar';
import ListSubheader from '@mui/material/ListSubheader';
import IconButton from '@mui/material/IconButton';
import { Info } from '@mui/icons-material';
import { InputLabel, RadioGroup, FormControlLabel, Radio, FormControl } from '@mui/material';
import CreatableSelect from 'react-select/creatable';
import { uploadAndSetTags } from '../uploadService.js'
import { storageRequestScope } from '../authConfig'
import Resizer from "react-image-file-resizer";
import EXIF from "exif-js";
import { getAlbumCollections } from "../photoService.js";
import styled from "styled-components";

const imageTypeRegex = /image\/(png|jpg|jpeg|JPEG|JPG|PNG)/gm;

const useStyles = makeStyles((theme) => ({
    root: {
        display: 'flex',
        flexWrap: 'wrap',
        justifyContent: 'space-around',
        overflow: 'hidden',
        backgroundColor: theme.palette.background.paper,
    },
    imageUpload: {
        margin: "auto",
        textAlign: "center",
        alignItems: "center",
        justifyContent: "center"
    },
    imageList: {
        display: 'flex',
        overflow: 'hidden',
        flexWrap: 'wrap'
    },
    inputs: {
        display: 'block',
        flexWrap: 'wrap',
        justifyContent: 'space-around'
    },
    imageItem: {
        height: 480,
        width: 640,
        "&:hover imageTitleBar": {
            display: 'inline-block'
        }
    },
    imageTitleBar: {
        display: 'inline-block',
        "&:hover": {

        }
    },
    icon: {
        color: 'rgba(255, 255, 255, 0.54)'
    },
    input: {
        accentColor: 'white',
        color: 'white',
        border: 'white'
    },
    imageInfo: {
        color: 'white',
        accentColor: 'white',
        border: 'white',
        backgroundColor: 'white'
    }
}));

const styledRadio = styled.input`
    color: 'white';
    accentColor: 'white';
    border: 'white';
    backgroundColor: 'white;
    width: 200px;
`;

export function ImageUpload() {
    const classes = useStyles();
    const [collectionAlbumData, setCollectionAlbumData] = useState([]);
    const [albumData, setAlbumData] = useState([]);
    const [currentAlbum, setCurrentAlbum] = useState([]);
    const [currentCollection, setCurrentCollection] = useState([]);
    const { instance, accounts } = useMsal();
    const [imageFiles, setImageFiles] = useState([]);
    const [images, setImages] = useState([]);
    const [isUploading, setIsUploading] = useState(false)
    const isMounted = useRef(true)
    const containerName = "uploads";
    const [isAlbumThumbnail, setIsAlbumThumbnail] = useState(['false']);

    const resizedFileHeight = 1800;
    const resizedFileWidth = 1200;
    const resizedImageQuality = 75;
    const resizedImageFormat = "JPEG";
    const resizedImageRotation = 0;

    const updateRadioThumbnailChanged = index => e => {
        console.log('index: ' + index);
        console.log('property name: ' + e.target.name);
        let newArr = []; // copying the old array
        // a deep copy is not needed as we are overriding the whole object below, and not setting a property of it. this does not mutate the state.
        newArr[index] = 'true';

        setIsAlbumThumbnail(newArr);
    }

    const resizeFile = (file, maxWidth, maxHeight, format, quality, rotation) =>
        new Promise((resolve) => {
            Resizer.imageFileResizer(
                file,
                maxWidth,
                maxHeight,
                format,
                quality,
                rotation,
                (uri) => {
                    resolve(uri);
                },
                "blob"
            );
        });

    const changeHandler = async (e) => {
        const { files } = e.target;
        const validImageFiles = [];
        for (let i = 0; i < files.length; i++) {
            const file = files[i];
            let exifData;
            if (file.type.match(imageTypeRegex)) {
                EXIF.getData(file, function () {
                    exifData = EXIF.getAllTags(file);
                });

                const blob = await resizeFile(file, resizedFileWidth, resizedFileHeight, resizedImageFormat, resizedImageQuality, resizedImageRotation);

                let resizedFile = new File([blob], file.name)
                let img = new Image();
                img.src = file;

                var fileData = {
                    file: resizedFile,
                    name: file.name,
                    exif: exifData,
                    isAlbumThumb: 'false'
                }

                validImageFiles.push(fileData);
            }
        }
        if (validImageFiles.length) {
            setImageFiles(validImageFiles);
            return;
        }
        alert("Selected images are not of valid type!");
    };

    // set isMounted to false when we unmount the component
    useEffect(() => {
        return () => {
            isMounted.current = false
        }
    }, [])

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

    const sendRequest = useCallback(async () => {
        if (isUploading) return
        setIsUploading(true)
        console.log("# imagefiles: " + imageFiles.length)
        console.log("album: " + currentAlbum)
        if (imageFiles.length > 0) {
            // add selected thumbnail data
            console.log("# imagefiles: " + imageFiles.length)
            for (let i = 0; i < imageFiles.length; i++) {
                imageFiles[i].isAlbumThumb = isAlbumThumbnail[i] === 'true' ? 'true' : 'false'
            }

            instance.acquireTokenSilent({
                ...storageRequestScope,
                account: accounts[0]
            }).then((response) => {
                uploadAndSetTags(containerName, imageFiles, currentCollection, currentAlbum)
                    .catch(error => console.log("error: " + error));
            })
        } else {
            return null
        }
        if (isMounted.current) // only update if we are still mounted
            setIsUploading(false)
    }, [isUploading, imageFiles, accounts, currentAlbum, currentCollection, containerName, instance, isAlbumThumbnail])

    useEffect(() => {
        const images = [], fileReaders = [];
        let isCancel = false;
        if (imageFiles.length) {
            imageFiles.forEach((file) => {
                const fileReader = new FileReader();
                fileReaders.push(fileReader);
                fileReader.onload = (e) => {
                    const image = new Image();
                    image.src = e.target.result;
                    //console.log(`image.width: ${image.width}`);
                    //console.log(`image.height: ${image.height}`);
                    const { result } = e.target;
                    //let fileObj = { "blob": result, "width": image.width, "height": image.height, "name": file.name }
                    let fileObj = { "blob": result, "name": file.name }
                    console.log(fileObj);
                    if (result) {
                        //images.push(result)
                        images.push(fileObj);
                    }
                    if (images.length === imageFiles.length && !isCancel) {
                        setImages(images);
                    }
                }
                fileReader.readAsDataURL(file.file);
            })
        };
        return () => {
            isCancel = true;
            fileReaders.forEach(fileReader => {
                if (fileReader.readyState === 1) {
                    fileReader.abort()
                }
            })
        }
    }, [imageFiles]);
    return (
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
    );
}