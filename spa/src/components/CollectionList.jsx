import React, { useState, useEffect } from 'react'
import { getAlbumCollections } from '../photoService';
import { List, ListItem, ListItemText } from '@mui/material';
import AlbumList from './AlbumList';
import { Route, Routes, Link } from 'react-router-dom';

export default function CollectionList() {
  const [collectionData, setCollectionData] = useState([])

  useEffect(() => {
    getAlbumCollections("uploads")
      .catch(error => console.log(error), (data = '[]') => console.log(data))
      .then((data) => setCollectionData(data));
  }, [])

  return (
    <>
      {
        Array.from(collectionData.keys()).map(function (collection, i) {
          return <>
            <p>
              <Link
                to={{
                  pathname: "/collections/albums",
                }}
                state={collectionData.get(collection)}
                key={i}
              >
                {collection}
              </Link>
            </p>
          </>
        })
      }
    </>
  )
}