import React from 'react'
import { useParams, useLocation, Link } from "react-router-dom"

export default function AlbumList() {
  let location = useLocation();
  let albums = location.state?.albums;
  let collection = location.state?.collection;
  {console.log("albums: " + albums)};

  return (
    <>
    Collection: {collection} 
    {console.log("collection: " + collection)}

      {albums.map((album, i) => {
            console.log("album: " + album);
        return <p key={album}>
          <Link
            to={{
              pathname: `/${collection}/${album}`
            }}
            state={{ collection: collection, album: album }}
            key={album}
          >
            {album}
          </Link>
        </p>
      })}
    </>
  )
}