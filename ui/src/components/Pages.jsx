import { Route, Routes } from 'react-router-dom';
import CollectionList from './CollectionList';
import AlbumList from './AlbumList';
import PhotoList from './PhotoList';
import { ImageUpload, UploadBlob } from './UploadBlob';

function Pages() {
    return (
        <Routes>
            <Route path="/" element={<CollectionList />} />
            <Route path="/:collection" element={<AlbumList />} />
            <Route path="/:collection/:album" element={<PhotoList />} />
            <Route path="/uploads" element={<UploadBlob />} />
        </Routes>
    );
}

export default Pages;
