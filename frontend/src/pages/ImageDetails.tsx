import {EuiButton, EuiFlexGrid, EuiPageTemplate, EuiText} from '@elastic/eui';
import {MutableRefObject, Ref, useRef, useState} from "react";
import {EuiFieldText, EuiFlexGroup, EuiFlexItem} from '@elastic/eui';
import {EuiInnerText} from "@elastic/eui";
import {useLocation, useNavigate} from "react-router";

function CreateImageDetailsPage() {
  const location = useLocation()
  const navigate = useNavigate()

  const buttonStyle = {
    marginRight: "25px",
    marginTop: "25px"
  }

  return(<>
    <EuiPageTemplate.Header pageTitle="Image Details" />

    <EuiFlexGroup>
      <EuiFlexItem grow={6}>
        <EuiPageTemplate.Section>
          <EuiText textAlign={"center"}>
            <strong>ID:</strong>
          </EuiText>
          <EuiText textAlign={"center"}>
            {location.state.item.id}
          </EuiText>
        </EuiPageTemplate.Section>

        <EuiPageTemplate.Section>
          <EuiText textAlign={"center"}>
            <strong>Name:</strong>
          </EuiText>
          <EuiText textAlign={"center"}>
            {location.state.item.name}
          </EuiText>
        </EuiPageTemplate.Section>

        <EuiPageTemplate.Section>
          <EuiText textAlign={"center"}>
            <strong>Operating System:</strong>
          </EuiText>
          <EuiText textAlign={"center"}>
            {location.state.item.operatingSystem}
          </EuiText>
        </EuiPageTemplate.Section>

        <EuiPageTemplate.Section>
          <EuiText textAlign={"center"}>
            <strong>Date Created:</strong>
          </EuiText>
          <EuiText textAlign={"center"}>
            {location.state.item.date.toString()}
          </EuiText>
        </EuiPageTemplate.Section>

        <EuiPageTemplate.Section>
          <EuiText textAlign={"center"}>
            <strong>Size:</strong>
          </EuiText>
          <EuiText textAlign={"center"}>
            {location.state.item.size}
          </EuiText>
        </EuiPageTemplate.Section>
      </EuiFlexItem>

      <EuiFlexItem>
        <EuiFlexGroup direction={"column"}>
          <EuiFlexItem grow={false}>
            <EuiButton 
            style={buttonStyle}
            onClick={() => {
              navigate('/createImage', {state:{item:location.state.item}})
            }}
            >
              Derive New Image
              </EuiButton>
          </EuiFlexItem>
          <EuiFlexItem grow={false}>
            <EuiButton style={buttonStyle}>Delete Image</EuiButton>
          </EuiFlexItem>
        </EuiFlexGroup>
      </EuiFlexItem>
    </EuiFlexGroup>

  </>)
}


export default CreateImageDetailsPage;