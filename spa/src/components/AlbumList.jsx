import React from 'react'
import { useParams, useLocation, Link } from "react-router-dom"

export default function AlbumList() {
  let location = useLocation();
  return (
    <>
      {console.log(location.state)}
      {location.state.map((album, i) => {
        return <p><Link key={i}>{album}</Link></p>
      })}
    </>
  )
}