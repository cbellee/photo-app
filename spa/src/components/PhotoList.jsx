import React, { useState, useEffect } from 'react'
import { getAlbumPhotos } from '../photoService';
import { useMsal } from "@azure/msal-react";
import { makeStyles, withTheme } from '@mui/styles';
import ImageList from '@mui/material/ImageList';
import ImageListItem from '@mui/material/ImageListItem';
import ImageListItemBar from '@mui/material/ImageListItemBar';
import ListSubheader from '@mui/material/ListSubheader';
import IconButton from '@mui/material/IconButton';
import { Info, NoEncryption } from '@mui/icons-material';

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

export default function PhotoList() {
  const classes = useStyles();
  const { instance, accounts } = useMsal();
  const [photoData, setPhotoData] = useState([])

  useEffect(() => {
    getAlbumPhotos()
      .catch(error => console.log(error), (data = '[]') => console.log(data))
      .then((data) => setPhotoData(data));
  }, [])

  return (
    <div className={classes.root}>
      <ImageList cellHeight={300} spacing={30} className={classes.imageList}>
        <ImageListItem key="Subheader" cols={2} style={{ height: 'auto' }}>
          <ListSubheader component="div"></ListSubheader>
        </ImageListItem>
        {photoData.map((photo) => (
          <ImageListItem key={photo.value.url} className={classes.imageItem}>
            <img src={photo.value.url} alt={photo.value.name} />
            <ImageListItemBar
              className={classes.imageTitleBar}
              title={photo.value.name}
              actionIcon={
                <IconButton aria-label={`info about ${photo.value.name}`} className={classes.icon}>
                  <Info />
                </IconButton>
              }
            />
          </ImageListItem>
        ))}
      </ImageList>
    </div>
  );
}