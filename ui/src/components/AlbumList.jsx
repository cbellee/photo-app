import { React, useState, useEffect } from 'react'
import { getAlbumPhoto } from '../services/photoService';
import { useLocation, Link } from "react-router-dom"
import { Typography } from '@mui/material';

export default function AlbumList() {
  const [albumPhotoData, setAlbumPhotoData] = useState(new Map())
  const updateMap = (key, value) => {
    setAlbumPhotoData(map => new Map(map.set(key, value)));
  }

  let location = useLocation();
  let albums = location.state?.albums;
  let collection = location.state?.collection;

  useEffect(() => {
    albums.map(function (album) {
      getAlbumPhoto("thumbs", collection, album)
        .catch(error => console.log(error), (data = '[]') => console.log(data))
        .then(data => updateMap(album, data));
    })
  }, [albums, collection]);

  return (
    <>
      <Typography>
        Albums
      </Typography>
      <Link
        to={{
          pathname: `/`
        }}
      >
        {`<< collections/${collection}`}
      </Link>
      {albums.map((album, i) => {
        return <p key={album}>
          <Link
            to={{
              pathname: `/${collection}/${album}`
            }}
            state={{ collection: collection, album: album, albums: albums }}
            key={album}
          >
            <li>{album}</li>
            <img src={albumPhotoData.get(album)} alt={album} />
          </Link>
        </p>
      })}
    </>
  )
}