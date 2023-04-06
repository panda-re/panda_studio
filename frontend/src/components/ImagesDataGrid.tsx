import {
  EuiBasicTable,
  EuiBasicTableColumn,
  EuiButton,
  EuiButtonIcon,
  EuiFieldText,
  EuiFilePicker,
  EuiFlexGroup,
  EuiFlexItem,
  EuiModal,
  EuiModalBody,
  EuiModalFooter,
  EuiModalHeader,
  EuiModalHeaderTitle,
  EuiOverlayMask,
  EuiSearchBar,
  EuiSearchBarOnChangeArgs,
  EuiSpacer,
  EuiText,
  RIGHT_ALIGNMENT,
  useGeneratedHtmlId
} from '@elastic/eui';
import {getItemId} from '@elastic/eui/src/components/basic_table/basic_table';
import {useQueryClient} from '@tanstack/react-query';
import axios, {AxiosRequestConfig} from 'axios';
import prettyBytes from 'pretty-bytes';
import React, {useEffect, useState} from 'react';
import {useLocation, useNavigate} from 'react-router-dom';
import {
  CreateImageFileFromUrlRequest,
  CreateImageFileRequest,
  CreateImageRequest,
  findAllImages,
  Image,
  ImageFile,
  ImageFileType,
  PandaConfig,
  updateImage,
  useCreateImage,
  useCreateImageFile, useCreateImageFileFromUrl,
  useDeleteImageById,
  useFindAllImages,
  useUpdateImage
} from '../api';

function ImagesDataGrid() {
  const navigate = useNavigate();
  const location = useLocation();
  const {isLoading, error, data} = useFindAllImages();
  const queryClient = useQueryClient();
  const deleteFunction = useDeleteImageById({mutation: {onSuccess: () => queryClient.invalidateQueries()}});
  const updateFn = useUpdateImage({mutation: {onSuccess: () => queryClient.invalidateQueries()}});
  const createFileFromUrl = useCreateImageFileFromUrl({mutation: {onSuccess() {
    setIsLoadingVisible(false);
    queryClient.invalidateQueries();
  }}})

  // File picker constants
  const createFileFn = useCreateImageFile({
    mutation: {
      onSuccess(data, variables, context) {
        setIsLoadingVisible(false);
        queryClient.invalidateQueries();
      }
    }
  })
  const filePickerId = useGeneratedHtmlId({prefix: 'filePicker'});
  const [files, setFiles] = useState(new Array<File>);

  const onFileChange = (files: FileList | null) => {
    setFiles(files!.length > 0 ? Array.from(files!) : []);
  };

  ///////// Modal Constants ///////////////////
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [modalName, setModalName] = useState("");
  const [modalDesc, setModalDesc] = useState("");
  const [modalArch, setModalArch] = useState("");
  const [modalOs, setModalOs] = useState("");
  const [modalPrompt, setModalPrompt] = useState("");
  const [modalCdrom, setModalCdrom] = useState("");
  const [modalSnapshot, setModalSnapshot] = useState("");
  const [modalMemory, setModalMemory] = useState("");
  const [modalExtraArgs, setModalExtraArgs] = useState("");
  const [url, setModalUrl] = useState("");

  const [isLoadingVisible, setIsLoadingVisible] = useState(false);

  const closeModal = () => {
    setModalName("");
    setModalDesc("");
    setModalArch("");
    setModalOs("");
    setModalPrompt("");
    setModalCdrom("");
    setModalSnapshot("");
    setModalMemory("");
    setModalExtraArgs("");
    setIsModalVisible(false)
    setModalUrl("")
  };
  const showModal = () => {
    setIsModalVisible(true);
  }


  /////////// Endpoint Functions //////////////
  const deleteImage = ({itemId}: { itemId: string }) => {
    deleteFunction.mutate({imageId: itemId});
  }

  const updateImage = ({image}: { image: Image }) => {
    if (image.id == null) {
      return;
    }
    const conf: PandaConfig = {
      qcow_file_name: image.config?.qcow_file_name,
      arch: modalArch,
      os: modalOs,
      prompt: modalPrompt,
      cdrom: modalCdrom,
      snapshot: modalSnapshot,
      memory: modalMemory,
      extra_args: modalExtraArgs,
    }
    const req: CreateImageRequest = {
      name: image.name,
      description: image.description,
      config: conf,
    };
    updateFn.mutate({data: req, imageId: image.id});
  }

  function deleteActionPress(event: React.MouseEvent, item: Image) {
    deleteImage({itemId: item.id!})
    event.stopPropagation();
  }

  function getFileTypeFromString(fileTypeAsString: string, imageId: string): ImageFileType | undefined {
    switch (fileTypeAsString) {
      case "qcow2": {
        return ImageFileType.qcow2;
      }
      case "dtb": {
        return ImageFileType.dtb;
      }
      case "kernel": {
        return ImageFileType.kernel;
      }
      default: {
        alert("Invalid File Type");
        deleteImage({itemId: imageId ?? ""})
        return;
      }
    }
  }


  function createImageFileFromUrl(image: Image) {
    const urlAsArray = url.split("/")
    const fileName = urlAsArray[urlAsArray.length - 1]
    const fileTypeAsArray = fileName.split(".")
    const fileTypeAsString = fileTypeAsArray[fileTypeAsArray.length - 1]
    const fileType = getFileTypeFromString(fileTypeAsString, image.id!)

    if (fileType == undefined) {
      return;
    }

    const fileFromUrlReq: CreateImageFileFromUrlRequest = {
      file_name: fileName,
      file_type: fileType,
      url: url
    }

    createFileFromUrl.mutate({data: fileFromUrlReq, imageId: image.id!}, {onError() {
      deleteImage({itemId: image.id!});
      setIsLoadingVisible(false);
      alert("Received an invalid URL");
    }})

    setIsLoadingVisible(true);
    closeModal();

  }

  function createFiles(image: Image) {
    var splitArray = files[0].name.split(".");
    var fileTypeString = splitArray[splitArray.length - 1];
    var fileType = getFileTypeFromString(fileTypeString, image.id!)

    if (fileType == undefined) {
      return;
    }

    const fileReq: CreateImageFileRequest = {
      file_name: files[0].name,
      file_type: fileType,
      file: files[0],
    }
    createFileFn.mutate({data: fileReq, imageId: image.id!})
    setIsLoadingVisible(true);
    closeModal();
  }

  const createFn = useCreateImage({
    mutation: {
      onSuccess(data, variables, context) {
        if (url == "") {
          createFiles(data)
        } else {
          createImageFileFromUrl(data)
        }
      },
    }
  })

  function createFile() {
    if (modalName == "" || modalArch == "" || modalOs == "" || modalPrompt == "" || modalMemory == "") {
      alert("Please fill out all required fields")
      return;
    }

    var fileName
    if (url == "") {
      fileName = files[0].name
    } else {
      const urlAsArray = url.split("/")
      const urlFileName = urlAsArray[urlAsArray.length - 1]
      fileName = urlFileName
    }

    const conf: PandaConfig = {
      qcow_file_name: fileName,
      arch: modalArch,
      os: modalOs,
      prompt: modalPrompt,
      cdrom: modalCdrom,
      snapshot: modalSnapshot,
      memory: modalMemory,
      extra_args: modalExtraArgs,
    }
    const req: CreateImageRequest = {
      name: modalName,
      description: modalDesc,
      config: conf,
    };
    createFn.mutate({data: req})
  }

  //////// UI Functions ///////////
  useEffect(() => {
    if (location.state) {
      if (location.state.image) {
        updateImage({image: location.state.image});
      } else {
        deleteImage({itemId: location.state.imageId});
      }
      window.history.replaceState({}, document.title)
    }
  }, []);

  function LoadingModal() {
    return <EuiOverlayMask>
      <EuiModal onClose={closeModal}>
        <EuiModalHeader>
          <EuiModalHeaderTitle>Uploading Image</EuiModalHeaderTitle>
        </EuiModalHeader>
        <EuiModalBody>
          <EuiText>
            Loading...
          </EuiText>
        </EuiModalBody>
      </EuiModal>
    </EuiOverlayMask>
  }

  function CreateModal() {
    return <EuiOverlayMask>
      <EuiModal onClose={closeModal}>
        <EuiModalHeader>
          <EuiModalHeaderTitle>Upload New Image</EuiModalHeaderTitle>
        </EuiModalHeader>
        <EuiModalBody>
          <EuiFieldText
            placeholder="Enter Name (required)"
            isInvalid={modalName == ""}
            name="imageName"
            onChange={(e) => {
              setModalName(e.target.value);
            }}/>
          <EuiFieldText
            placeholder="Enter New Description"
            name="imageDesc"
            onChange={(e) => {
              setModalDesc(e.target.value);
            }}/>
          <EuiFieldText
            placeholder="Enter image Architecture (required)"
            isInvalid={modalArch == ""}
            name="pandaConfigArch"
            onChange={(e) => {
              setModalArch(e.target.value);
            }}/>
          <EuiFieldText
            placeholder="Enter image OS (required)"
            isInvalid={modalOs == ""}
            name="pandaConfigOs"
            onChange={(e) => {
              setModalOs(e.target.value);
            }}/>
          <EuiFieldText
            placeholder="Enter prompt (required)"
            isInvalid={modalPrompt == ""}
            name="pandaConfigPrompt"
            onChange={(e) => {
              setModalPrompt(e.target.value);
            }}/>
          <EuiFieldText
            placeholder="Enter Cdrom"
            name="pandaConfigCdrom"
            onChange={(e) => {
              setModalCdrom(e.target.value);
            }}/>
          <EuiFieldText
            placeholder="Enter Snapshot"
            name="pandaConfigSnapshot"
            onChange={(e) => {
              setModalSnapshot(e.target.value);
            }}/>
          <EuiFieldText
            placeholder="Enter memory amount (required)"
            isInvalid={modalMemory == ""}
            name="pandaConfigMemory"
            onChange={(e) => {
              setModalMemory(e.target.value);
            }}/>
          <EuiFieldText
            placeholder="Enter Extra args"
            name="pandaConfigExtraArgs"
            onChange={(e) => {
              setModalExtraArgs(e.target.value);
            }}/>
          <EuiFilePicker
            id={filePickerId}
            initialPromptText="Select or drag and drop multiple files"
            onChange={onFileChange}
            aria-label="Use aria labels when no actual label is in use"
          />
          <EuiText>Alternatively, use a URL to a valid image file:</EuiText>
          <EuiFieldText placeholder={"Enter an image URL"} onChange={(e) => {
            setModalUrl(e.target.value);
          }}/>

        </EuiModalBody>
        <EuiModalFooter>
          <EuiButton onClick={closeModal} fill>Close</EuiButton>
          <EuiButton
            onClick={() => {
              createFile();
            }}
            fill>
            Submit</EuiButton>
        </EuiModalFooter>
      </EuiModal>
    </EuiOverlayMask>
  }

  const tableColumns: EuiBasicTableColumn<Image>[] = [
    {
      field: 'id',
      name: 'Id',
    },
    {
      field: 'name',
      name: 'File Name',
    },
    {
      field: 'files',
      name: 'Size',
      render: (value: ImageFile[]) => {
        var size = 0;
        for (var f of value) {
          size += (f.size != null) ? +f.size : 0;
        }
        return prettyBytes(size, {maximumFractionDigits: 2});
      },
    },
    {
      align: RIGHT_ALIGNMENT,
      name: 'Delete',
      render: (item: Image) => {
        return (
          <EuiButtonIcon
            onClick={(event: React.MouseEvent) => {
              deleteActionPress(event, item)
            }}
            iconType={"trash"}
          />
        );
      },
    },
  ]

  const getRowProps = (item: Image) => {
    const id = item.id;
    return {
      'data-test-subj': `image-row-${id}`,
      onClick: () => {
        navigate('/imageDetails', {state: {item: item}})
      },
    }
  }

  //////////// Search Bar Items //////////////
  const initialQuery = EuiSearchBar.Query.MATCH_ALL;

  const [query, setQuery] = useState(initialQuery);

  const onChange = (args: EuiSearchBarOnChangeArgs) => {
    setQuery(args.query ?? initialQuery);
  };

  const queriedItems = EuiSearchBar.Query.execute(query, data ?? []);

  ////////// UI Element //////////
  return (<>
    <EuiFlexGroup justifyContent='spaceBetween'>
      <EuiFlexItem grow={false} style={{minWidth: 300}}>
        <EuiSearchBar
          box={{
            incremental: true,
          }}
          defaultQuery={initialQuery}
          onChange={onChange}/>
      </EuiFlexItem>
      <EuiFlexItem grow={false}>
        <EuiButton onClick={showModal} iconType={'plusInCircle'}>Upload Base Image</EuiButton>
      </EuiFlexItem>
    </EuiFlexGroup>
    <EuiSpacer></EuiSpacer>
    {isLoading && <div>Loading...</div> ||
    <EuiBasicTable
      tableCaption="Images"
      items={queriedItems ?? []}
      rowHeader="firstName"
      columns={tableColumns}
      rowProps={getRowProps}
    />
    }
    {(isModalVisible) ? (CreateModal()) : null}
    {(isLoadingVisible) ? (LoadingModal()) : null}
  </>)
}

export default ImagesDataGrid;