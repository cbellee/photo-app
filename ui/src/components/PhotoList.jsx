import React, { useState, useEffect } from 'react'
import { useLocation, Link } from 'react-router-dom'
import { getAlbumCollections } from '../photoService';
import { useMsal } from "@azure/msal-react";
import { makeStyles, withTheme } from '@mui/styles';
import { Modal, Fade, Backdrop, linearProgressClasses } from '@mui/material';
import ImageList from '@mui/material/ImageList';
import ImageListItem from '@mui/material/ImageListItem';
import ImageListItemBar from '@mui/material/ImageListItemBar';
import ListSubheader from '@mui/material/ListSubheader';
import IconButton from '@mui/material/IconButton';
import { Info, NoEncryption } from '@mui/icons-material';
import { getBlobsByTags } from '../photoService';
import styled, { keyframes, css } from 'styled-components'

const useStyles = makeStyles((theme) => ({
  modal: {
    maxWidth: '100%',
    margin: 'auto',
    position: "absolute",
    width: 'auto',
    backgroundColor: theme.palette.background.paper,
    boxShadow: theme.shadows[10],
    padding: theme.spacing(3, 3, 3, 3),
    borderRadius: 3,
    "&:focus-visible": {
      outline: 'white auto 0px',
    }
  },
  modalOverlay: {
    backgroundColor: 'black',
  },
  image: {
    maxHeight: '40vw'
  },
  root: {
    display: 'flex',
    flexWrap: 'wrap',
    justifyContent: 'space-around',
    overflow: 'hidden',
    backgroundColor: theme.palette.background.paper,
  },
  imageList: {
    background: 'white'
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

function getModalStyle() {
  const top = 50;
  const left = 50;

  return {
    top: `${top}%`,
    left: `${left}%`,
    transform: `translate(-${top}%, -${left}%)`
  };
}

export default function PhotoList() {
  const classes = useStyles();
  const { instance, accounts } = useMsal();
  const [photoData, setPhotoData] = useState([]);
  let location = useLocation();
  const [open, setOpen] = useState(false);
  let [currentIdx, setCurrentIdx] = useState(null);
  const [modalStyle] = useState(getModalStyle);

  const handleOpen = (idx) => {
    setCurrentIdx(idx);
    setOpen(true);
  };

  const handleClose = () => {
    setCurrentIdx(null);
    setOpen(false);
  };

  const handleScroll = (idx) => {
    if (currentIdx >= photoData.length - 1) {
      setCurrentIdx(0)
    } else {
      setCurrentIdx(++idx);
    }
  };

  useEffect(() => {
    let collection = location.state.collection;
    let album = location.state.album;

    getBlobsByTags("uploads", collection, album, true)
      .catch(error => console.log(error), (data = '[]') => console.log("data: " + data))
      .then((data) => setPhotoData(data));
  }, [])

  const body = (
    <div style={modalStyle} className={classes.modal}>
      {photoData[currentIdx] && (
        <img className={classes.image} src={photoData[currentIdx].tags.url} alt={photoData[currentIdx].name} onClick={() => handleScroll(currentIdx)} />
      )}
    </div>
  );

  return (
    <div className={classes.root}>
      <ImageList cellHeight={300} spacing={30} className={classes.imageList}>
        <ImageListItem key="Subheader" cols={4} style={{ height: 'auto' }}>
          <ListSubheader component="div"></ListSubheader>
        </ImageListItem>
        {photoData.map((photo, idx) => (
          <ImageListItem key={photo.tags.thumbUrl} className={classes.imageItem}>
            <img
              src={photo.tags.thumbUrl}
              alt={photo.name}
              onClick={() => handleOpen(idx)}
            />
            <ImageListItemBar
              className={classes.imageTitleBar}
              title={photo.name}
              actionIcon={
                <IconButton aria-label={`info about ${photo.name}`} className={classes.icon}>
                  <Info />
                </IconButton>
              }
            />
          </ImageListItem>
        ))}
      </ImageList>
      <Modal
        open={open}
        className={classes.modalOverlay}
        onClose={handleClose}
        onClick={() => handleScroll(currentIdx)}
        center
      >
        <Fade in={open}>
          {body}
        </Fade>
      </Modal>
    </div >
  );
}