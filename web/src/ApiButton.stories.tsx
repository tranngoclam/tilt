import { Story } from "@storybook/react"
import React from "react"
import { MemoryRouter } from "react-router"
import { ApiButton } from "./ApiButton"

type UIButton = Proto.v1alpha1UIButton

export default {
  title: "New UI/Shared/ApiButton",
  decorators: [
    (Story: any) => (
      <MemoryRouter initialEntries={["/"]}>
        {/* required for MUI <Icon> */}
        <link
          rel="stylesheet"
          href="https://fonts.googleapis.com/icon?family=Material+Icons"
        />
        <div style={{ margin: "-1rem" }}>
          <Story />
        </div>
      </MemoryRouter>
    ),
  ],
}

function newButton(icon?: string, text?: string): UIButton {
  return {
    metadata: {
      name: "button",
    },
    spec: {
      iconName: icon,
      text: text,
    },
  }
}

type StoryProps = {
  text?: string
  iconName?: string
  showText: boolean
}
const Template: Story<StoryProps> = (args) => (
  <ApiButton
    button={newButton(args.iconName, args.text)}
    showText={args.showText}
  />
)

export const IconAndText = Template.bind({})
IconAndText.args = {
  text: "yum!",
  iconName: "cake",
  showText: true,
}

export const OnlyIcon = Template.bind({})
OnlyIcon.args = {
  iconName: "cake",
}

export const OnlyText = Template.bind({})
OnlyText.args = {
  text: "yum!",
  showText: true,
}
