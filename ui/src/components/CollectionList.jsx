import React, { useState, useEffect } from 'react'
import { getAlbumCollections, getCollectionPhoto } from '../services/photoService';
import { Link } from 'react-router-dom';
import { Typography } from '@mui/material';

export default function CollectionList() {
  const [collectionData, setCollectionData] = useState([])
  const [collectionPhotoData, setCollectionPhotoData] = useState(new Map())
  const updateMap = (key, value) => {
    setCollectionPhotoData(map => new Map(map.set(key, value)));
  }

  useEffect(() => {
    getAlbumCollections("uploads")
      .catch(error => console.log(error), (data = '[]') => console.log(data))
      .then((data) => setCollectionData(data));
  }, [])

  useEffect(() => {
    Array.from(collectionData.keys()).map(function (collection) {
      getCollectionPhoto("thumbs", collection)
        .catch(error => console.log(error), (data = '[]') => console.log(data))
        .then(data => updateMap(collection, data));
    })
  }, [collectionData])

  return (
    <>
      <Typography>
        Collections
      </Typography>
      {
        Array.from(collectionData.keys()).map(function (collection) {
          return <div key={collection}>
            <Link
              src={collectionPhotoData.get(collection)}
              to={{
                pathname: `/${collection}`,
              }}
              state={{ albums: collectionData.get(collection), collection: collection }}
              key={collection}
            >
              <ul>
                <li>{collection}</li>
              </ul>
              <img key={collection} src={collectionPhotoData.get(collection)} alt={collection} />
            </Link>
          </div>
        })
      }
    </>
  )
}