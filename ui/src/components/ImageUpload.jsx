import { useEffect, useState, useRef, useCallback } from "react";
import ImageList from '@mui/material/ImageList';
import { makeStyles } from '@mui/styles';
import { useMsal } from "@azure/msal-react";
import ImageListItem from '@mui/material/ImageListItem';
import ImageListItemBar from '@mui/material/ImageListItemBar';
import Select from 'react-select'
import ListSubheader from '@mui/material/ListSubheader';
import IconButton from '@mui/material/IconButton';
import { Info } from '@mui/icons-material';
import { InputLabel, Link } from '@mui/material';
import { DropdownButton, Dropdown } from "react-bootstrap/esm";
import CreatableSelect from 'react-select/creatable';
import { uploadAndSetTags } from '../uploadService.js'
import { storageRequestScope } from '../authConfig'
import Resizer from "react-image-file-resizer";
import EXIF from "exif-js";
import { Autocomplete, Container, TextField, Box, Typography } from "@mui/material";
import { getAlbumCollections } from "../photoService.js";
import { Form } from "react-router-dom";

const imageTypeRegex = /image\/(png|jpg|jpeg)/gm;

const useStyles = makeStyles((theme) => ({
    root: {
        display: 'flex',
        flexWrap: 'wrap',
        justifyContent: 'space-around',
        overflow: 'hidden',
        backgroundColor: theme.palette.background.paper,
    },
    imageList: {
        background: 'red'
    },
    imageItem: {
        height: 480,
        width: 640,
        "&:hover imageTitleBar": {
            display: 'block'
        }
    },
    imageTitleBar: {
        display: 'none !important',
        "&:hover": {
            background: 'white',
        }
    },
    icon: {
        color: 'rgba(255, 255, 255, 0.54)',
    },
}));

export function ImageUpload() {
    const classes = useStyles();
    const [collectionAlbumData, setCollectionAlbumData] = useState([]);
    const [collectionData, setCollectionData] = useState([]);
    const [albumData, setAlbumData] = useState([]);
    const [currentAlbum, setCurrentAlbum] = useState([]);
    const [currentCollection, setCurrentCollection] = useState([]);
    const { instance, accounts } = useMsal();
    const [imageFiles, setImageFiles] = useState([]);
    const [images, setImages] = useState([]);
    const [isUploading, setIsUploading] = useState(false)
    const isMounted = useRef(true)
    const containerName = "uploads";

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
                    console.log("all metadata: " + JSON.stringify(exifData));
                });

                const blob = await resizeFile(file, 1800, 1200, "JPEG", 75, 0);
                var resizedFile = new File([blob], file.name)

                var fileAndExif = {
                    file: resizedFile,
                    exif: JSON.stringify(exifData)
                }

                validImageFiles.push(fileAndExif);
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
    })

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
        setCurrentAlbum(event.value);
        console.log("current album: " + event.value)
    }

    const sendRequest = useCallback(async () => {
        if (isUploading) return
        setIsUploading(true)
        console.log("# imagefiles: " + imageFiles.length)
        if (imageFiles.length > 0) {
            console.log("# imagefiles: " + imageFiles.length)
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
    }, [isUploading, imageFiles, accounts])

    useEffect(() => {
        const images = [], fileReaders = [];
        let isCancel = false;
        if (imageFiles.length) {
            imageFiles.forEach((file) => {
                const fileReader = new FileReader();
                fileReaders.push(fileReader);
                fileReader.onload = (e) => {
                    const { result } = e.target;
                    if (result) {
                        images.push(result)
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
        <div className="ImageUpload">
            <form>
                <p>
                    <input
                        type="file"
                        id="file"
                        onChange={changeHandler}
                        accept="image/png, image/jpg, image/jpeg"
                        multiple
                    />
                </p>
                <label htmlFor="file">Upload images</label>
                {
                    imageFiles.length > 0 ?

                        <input
                            type="button"
                            id="upload"
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
            </form>
            <div className={classes.root}>
                {
                    images.length > 0 ?
                        <ImageList cellHeight={300} spacing={30} className={classes.imageList}>
                            <ImageListItem key="Subheader" cols={4} style={{ height: 'auto' }}>
                                <ListSubheader component="div"></ListSubheader>
                            </ImageListItem>
                            {images.map((image, idx) => {
                                    return <ImageListItem key={idx} className={classes.imageItem}>
                                        <img src={image} alt={image} />
                                        <ImageListItemBar
                                            className={classes.imageTitleBar}
                                            title={image}
                                            actionIcon={<IconButton aria-label={`info about ${image}`} className={classes.icon}>
                                                <Info />
                                            </IconButton>} />
                                </ImageListItem>
                            })}
            </ImageList> : null
                }
        </div>
        </div >
    );
}