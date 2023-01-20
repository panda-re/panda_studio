import {EuiButton, EuiPageTemplate, EuiText} from '@elastic/eui';
import {MutableRefObject, Ref, useRef, useState} from "react";
import {EuiFieldText, EuiFlexGroup, EuiFlexItem} from '@elastic/eui';
import {EuiInnerText} from "@elastic/eui";
import {useLocation} from "react-router";

function CreateImageDetailsPage() {
  const location = useLocation()
  console.log(location.state)

  return(<>
    <EuiPageTemplate.Header pageTitle="Image Details" />

    <EuiPageTemplate.Section>
      <EuiText>
        <strong>ID</strong>
      </EuiText>
      <EuiText textAlign={"center"}>
        {location.state.item.id}
      </EuiText>
    </EuiPageTemplate.Section>

    <EuiPageTemplate.Section>
      <EuiText>
        <strong>Name</strong>
      </EuiText>
      <EuiText textAlign={"center"}>
        {location.state.item.name}
      </EuiText>
    </EuiPageTemplate.Section>

    <EuiPageTemplate.Section>
      <EuiText>
        <strong>Operating System</strong>
      </EuiText>
      <EuiText textAlign={"center"}>
        {location.state.item.operatingSystem}
      </EuiText>
    </EuiPageTemplate.Section>

    <EuiPageTemplate.Section>
      <EuiText>
        <strong>Date Created</strong>
      </EuiText>
      <EuiText textAlign={"center"}>
        {location.state.item.date.toString()}
      </EuiText>
    </EuiPageTemplate.Section>

    <EuiPageTemplate.Section>
      <EuiText>
        <strong>Size</strong>
      </EuiText>
      <EuiText textAlign={"center"}>
        {location.state.item.size}
      </EuiText>
    </EuiPageTemplate.Section>

  </>)
}


export default CreateImageDetailsPage;