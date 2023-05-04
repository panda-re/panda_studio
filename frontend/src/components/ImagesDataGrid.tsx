import { EuiBasicTable, EuiBasicTableColumn, EuiButton, EuiButtonIcon, EuiConfirmModal, EuiFieldText, EuiFilePicker, EuiFlexGroup, EuiFlexItem, EuiModal, EuiModalBody, EuiModalFooter, EuiModalHeader, EuiModalHeaderTitle, EuiOverlayMask, EuiSearchBar, EuiSearchBarOnChangeArgs, EuiSelect, EuiSpacer, EuiText, RIGHT_ALIGNMENT, useGeneratedHtmlId } from '@elastic/eui';
import { useQueryClient } from '@tanstack/react-query';
import prettyBytes from 'pretty-bytes';
import React, {useEffect, useState} from 'react';
import {useLocation, useNavigate} from 'react-router-dom';
import {
  CreateImageFileFromUrlRequest,
  CreateImageFileRequest,
  CreateImageRequest,
  Image,
  ImageFile,
  ImageFileType,
  PandaConfig,
  useCreateImage,
  useCreateImageFile, useCreateImageFileFromUrl,
  useDeleteImageById,
  useFindAllImages,
  useUpdateImage
} from '../api';
import { archOptions } from './DefaultImageData';

function ImagesDataGrid() {
  const navigate = useNavigate();
  const location = useLocation();
  const {isLoading, isError, data} = useFindAllImages();
  const queryClient = useQueryClient();

  const deleteFunction = useDeleteImageById({
    mutation: {
      onSuccess: () => queryClient.invalidateQueries(),
      onError: (response) => alert("Error deleting Image:\n" + response)}});
  const updateFn = useUpdateImage({
    mutation: {
      onSuccess: () => queryClient.invalidateQueries(),
      onError: (response) => alert("Error updating image: \n" + response)}});
  const createFileFromUrl = useCreateImageFileFromUrl({
    mutation: {
      onSuccess() {
        setIsLoadingVisible(false);
        queryClient.invalidateQueries();
      },
      onError: (response) => alert("Error uploading image: \n" + response.error?.message)}});
      
   // File picker constants
   const createFileFn = useCreateImageFile({
    mutation: {
      onSuccess(data, variables, context) {
        setIsLoadingVisible(false);
        queryClient.invalidateQueries();
  }}})
   const filePickerId = useGeneratedHtmlId({ prefix: 'filePicker' });
   const [files, setFiles] = useState(new Array<File>);
 
   const onFileChange = (files: FileList | null) => {
     setFiles(files!.length > 0 ? Array.from(files!) : []);
   };

    const [archValue, setArchValue] = useState(archOptions[0].value);

    const basicSelectId = useGeneratedHtmlId({ prefix: 'basicSelect' });

    const onDropdownChange = (val: string) => {
      setArchValue(val);
    };
   
   ///////// Modal Constants ///////////////////
   const [isModalVisible, setIsModalVisible] = useState(false);
   const [modalName, setModalName] = useState("");
   const [modalDesc, setModalDesc] = useState("");
   const [modalOs, setModalOs] = useState("");
   const [modalPrompt, setModalPrompt] = useState("");
   const [modalCdrom, setModalCdrom] = useState("");
   const [modalSnapshot, setModalSnapshot] = useState("");
   const [modalMemory, setModalMemory] = useState("");
   const [modalExtraArgs, setModalExtraArgs] = useState("-display none");
   const [url, setModalUrl] = useState("");
 
   const [isLoadingVisible, setIsLoadingVisible] = useState(false);
   const [isConfirmVisible, setIsConfirmVisible] = useState(false);
   const [itemToDelete, setItemToDelete] = useState({})

   const closeModal = () => {
     setModalName("");
     setModalDesc("");
     setModalOs("");
     setModalPrompt("");
     setModalCdrom("");
     setModalSnapshot("");
     setModalMemory("");
     setModalExtraArgs("-display none");
     setIsModalVisible(false)
     setModalUrl("")
   };
   const showModal = () => {
     setIsModalVisible(true);
   }

   function showConfirmModal(event: React.MouseEvent, item: Image){
    setItemToDelete(item);
    setIsConfirmVisible(true);
    event.stopPropagation();
   }

   function closeConfirmModal(){
    setIsConfirmVisible(false);
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
      qcowfilename: image.config?.qcowfilename,
      arch: image.config?.arch,
      os: image.config?.os,
      prompt: image.config?.prompt,
      cdrom: image.config?.cdrom,
      snapshot: image.config?.snapshot,
      memory: image.config?.memory,
      extraargs: image.config?.extraargs,      
    }
    const req: CreateImageRequest = {
      name: image.name,
      description: image.description,
      config: conf,
    };
    updateFn.mutate({data: req, imageId: image.id});
  }

  function deleteActionPress(item: Image) {
    deleteImage({itemId: item.id!});
    setIsConfirmVisible(false);
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

  function createFile(){
    if(modalName=="" || modalSnapshot=="" || modalPrompt=="" || modalMemory==""){
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
      qcowfilename: fileName,
      arch: archValue,
      os: modalOs,
      prompt: modalPrompt,
      cdrom: modalCdrom,
      snapshot: modalSnapshot,
      memory: modalMemory,
      extraargs: modalExtraArgs,   
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

  function ConfirmModal(){
    return <EuiConfirmModal
        title="Are you sure you want to delete?"
        onCancel={closeConfirmModal}
        onConfirm={() => deleteActionPress(itemToDelete)}
        cancelButtonText="Cancel"
        confirmButtonText="Delete Image"
        buttonColor="danger"
        defaultFocusedButton="confirm"
      ></EuiConfirmModal>;
  }

  function LoadingModal() {
    return <EuiOverlayMask>
              <EuiModal onClose={()=>{}}>
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
                  <EuiFlexGroup>
                    {archOptions.map((element) => {
                      return <EuiFlexItem>
                                <EuiButton
                                  onClick={() => {
                                    setArchValue(element.value)
                                    setModalOs(element.defaultConfig.os);
                                    setModalPrompt(element.defaultConfig.prompt);
                                    setModalCdrom(element.defaultConfig.cdrom);
                                    setModalSnapshot(element.defaultConfig.snapshot);
                                    setModalMemory(element.defaultConfig.default_mem);
                                    setModalExtraArgs(element.defaultConfig.extra_args);
                                    setModalUrl(element.defaultConfig.url);
                                  }}
                                >{element.text}</EuiButton>
                                </EuiFlexItem>;
                    })}
                  </EuiFlexGroup>
                  <EuiSpacer size='m'></EuiSpacer>
                    <EuiFieldText 
                      placeholder="Enter Name (required)"
                      isInvalid={modalName == ""}
                      name="imageName"
                      value={modalName}
                      onChange={(e) => {
                        setModalName(e.target.value);
                      }}/>
                      <EuiFieldText 
                      placeholder="Enter New Description"  
                      name="imageDesc"
                      value={modalDesc}
                      onChange={(e) => {
                        setModalDesc(e.target.value);
                      }}/>
                      <EuiSelect
                        id={basicSelectId}
                        options={archOptions}
                        value={archValue}
                        onChange={(e) => {
                          onDropdownChange(e.target.value);
                        }}/>
                      <EuiFieldText 
                      placeholder="Enter image OS"
                      name="pandaConfigOs"
                      value={modalOs}
                      onChange={(e) => {
                        setModalOs(e.target.value);
                      }}/>
                      <EuiFieldText 
                      placeholder="Enter prompt (required)"
                      isInvalid={modalPrompt == ""}
                      name="pandaConfigPrompt"
                      value={modalPrompt}
                      onChange={(e) => {
                        setModalPrompt(e.target.value);
                      }}/>
                      <EuiFieldText 
                      placeholder="Enter Cdrom" 
                      name="pandaConfigCdrom" 
                      value={modalCdrom}
                      onChange={(e) => {
                        setModalCdrom(e.target.value);
                      }}/>
                      <EuiFieldText 
                      placeholder="Enter Snapshot (required)"  
                      name="pandaConfigSnapshot" 
                      value={modalSnapshot}
                      isInvalid={modalSnapshot == ""}
                      onChange={(e) => {
                        setModalSnapshot(e.target.value);
                      }}/>
                      <EuiFieldText 
                      placeholder="Enter memory amount (required)"
                      isInvalid={modalMemory == ""}
                      name="pandaConfigMemory" 
                      value={modalMemory}
                      onChange={(e) => {
                        setModalMemory(e.target.value);
                      }}/>
                      <EuiFieldText 
                      placeholder="Enter Extra args"  
                      name="pandaConfigExtraArgs"
                      value={modalExtraArgs}
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
                      <EuiFieldText 
                        placeholder={"Enter an image URL"} 
                        value={url}
                        onChange={(e) => {
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
              showConfirmModal(event, item);
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
    {(isError) ? (<div>Error...</div>)
    : ((isLoading) ? (<div>Loading...</div>) 
    : <EuiBasicTable
      tableCaption="Images"
      items={queriedItems ?? []}
      rowHeader="firstName"
      columns={tableColumns}
      rowProps={getRowProps}
    />)
  }
  {(isModalVisible) ? (CreateModal()) : null}
  {(isLoadingVisible) ? (LoadingModal()) : null}
  {(isConfirmVisible) ? (ConfirmModal()) : null}
  </>)
}

export default ImagesDataGrid;