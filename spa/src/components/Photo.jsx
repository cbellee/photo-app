export const Photo = ({ name, url, id }) => (
    <div>
      <img src={url} />
      <div>
        <p>{name}</p>
        <p>{url}</p>
        <p>{id}</p>
      </div>
    </div>
  );