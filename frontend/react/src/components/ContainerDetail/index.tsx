import styles from "./styles.module.css";

interface Props {
  handleClose: () => void;
  containerName: string;
  containerImage: string;
  containerStatus: string;
  containerStarted: string;
  containerPorts: string;
  containerId: string;
}

const ContainerDetail: React.FC<Props> = ({
  handleClose,
  containerName,
  containerImage,
  // containerStarted,
  containerStatus,
  containerPorts,
  containerId,
}) => {
  function handleDelete() {
    fetch("/deleteContainer", {
      method: "POST",
      body: JSON.stringify({ container_id: containerId }),
    }).then((res) => {
      if (res.status === 200) {
        alert("Container deleted");
      } else {
        alert("Error deleting container");
      }
      window.location.reload();
    });
  }

  function handleStop() {
    fetch("/stopContainer", {
      method: "POST",
      body: JSON.stringify({ container_id: containerId }),
    }).then((res) => {
      if (res.status === 200) {
        alert("Container stopped");
      } else {
        alert("Error stopping container");
      }
      window.location.reload();
    });
  }

  function handleEdit() {
    const newContainerName = prompt("Enter new container name");
    if (newContainerName) {
      fetch("/editContainer", {
        method: "POST",
        body: JSON.stringify({
          container_id: containerId,
          new_container_name: newContainerName,
        }),
      }).then((res) => {
        if (res.status === 200) {
          alert("Container edited");
        } else {
          alert("Error editing container");
        }
        window.location.reload();
      });
    }
  }

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h1 className={styles.title}>{containerName}</h1>
        <button onClick={handleClose} className={styles.closeBtn}>
          X
        </button>
      </div>
      <div className={styles.infoContainer}>
        <div className={styles.infoField}>
          <p className={styles.infoTitle}>ID:</p>
          <p>{containerId}</p>
        </div>
        <div className={styles.infoField}>
          <p className={styles.infoTitle}>Ports:</p>
          <p>{containerPorts}</p>
        </div>
        <div className={styles.infoField}>
          <p className={styles.infoTitle}>Status:</p>
          <p>{containerStatus}</p>
        </div>
        <div className={styles.infoField}>
          <p className={styles.infoTitle}>Image:</p>
          <p>{containerImage}</p>
        </div>
        {/* <p>Started: {containerStarted}</p> */}
      </div>
      <div className={styles.buttons}>
        <button
          onClick={() => {
            window.open(
              `${window.location.protocol}//${window.location.hostname}:${
                containerPorts.split(":")[1].split("->")[0]
              }`,
              "_blank"
            );
          }}
        >
          Open Service
        </button>
        <button onClick={handleEdit}>Edit Container</button>
        <button onClick={handleStop}>Stop Container</button>
        <button onClick={handleDelete}>Delete Container</button>
      </div>
    </div>
  );
};

export default ContainerDetail;
