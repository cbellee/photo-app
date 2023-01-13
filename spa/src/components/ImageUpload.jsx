import { useEffect, useState, useRef, useCallback } from "react";
import ImageList from '@mui/material/ImageList';
import { makeStyles } from '@mui/styles';
import { useMsal } from "@azure/msal-react";
import ImageListItem from '@mui/material/ImageListItem';
import ImageListItemBar from '@mui/material/ImageListItemBar';
import ListSubheader from '@mui/material/ListSubheader';
import IconButton from '@mui/material/IconButton';
import { Info } from '@mui/icons-material';
import { uploadAndSetTags } from '../uploadService.js'
import { storageRequestScope } from '../authConfig'
import Resizer from "react-image-file-resizer";
import EXIF from "exif-js";

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
                uploadAndSetTags(containerName, imageFiles, "athletics", "2022")
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
            </form>
            <div className={classes.root}>
                {
                    images.length > 0 ?
                        <ImageList cellHeight={300} spacing={30} className={classes.imageList}>
                            <ImageListItem key="Subheader" cols={2} style={{ height: 'auto' }}>
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
                                </ImageListItem>;
                            })}
                        </ImageList> : null
                }
            </div>
        </div>
    );
}