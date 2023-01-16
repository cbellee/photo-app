import React, { useState, useEffect } from 'react'
import { getAlbumCollections } from '../photoService';
import { Link } from 'react-router-dom';

export default function CollectionList() {
  const [collectionData, setCollectionData] = useState([])

  useEffect(() => {
    getAlbumCollections("uploads")
      .catch(error => console.log(error), (data = '[]') => console.log(data))
      .then((data) => setCollectionData(data));
  }, [])

  return (
    <>
    Collections
      {
        Array.from(collectionData.keys()).map(function (collection, i) {
          return <>
            <p>
              <Link
                to={{
                  pathname: `/${collection}`,
                }}
                state={{ albums: collectionData.get(collection), collection: collection }}
                key={collection}
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