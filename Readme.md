# ⚡ BashForm

**Create and share forms using SSH**. Users can participate in forms without installing any additional packages.

## Key Features

- **💻 Terminal-Based Forms**: Create forms and fill them out directly from the terminal.
- **🔑 Secure Authorization**: Users authenticate using their SSH keys, eliminating the need for passwords and ensuring a secure experience.
- **👶 No Installation Required**: Users respond to forms via SSH, eliminating the need for installing client software.
- **🫂 Easy Form Sharing**: Share forms using unique codes, enabling global participation.
- **🚀 Customizable Forms**: Create forms with any number of questions, tailored to your needs.

## How It Works

Bashform leverages SSH to provide a simple and secure interface for creating and interacting with forms. Here's how:

- **Form Creation**: Generate forms with a specified number of questions and a unique code.
- **Form Filling**: Respond to forms securely via SSH using the assigned code.
- **Form Responses**: Retrieve form responses, allowing for easy data collection.

<p align="center">
    <img height="400" src="/images/open.gif">
</p>

## Installation

There is no installation required for using Bashform. As long as you have SSH access, you can:

- Create forms
- Respond to forms
- Get form responses

## Usage

> [!NOTE]
> You need an SSH key to use bashform. If you don't have one (WHY?), you can generate one using the following command:
>
> ```bash
> ssh-keygen -t rsa -b 4096 -C "<your_email>"
> ```

### Fill Out a Form

To fill out a form, use the following command:

```bash
ssh -t bashform.me form <code>
# or
ssh -t bashform.me f <code>
```

Replace `<code>` with the unique code of the form you wish to fill out.

<p align="center">
    <img height="600" src="/images/form.gif">
</p>

### Create a New Form

To create a new form, use the command:

```bash
ssh -t bashform.me create <num_of_questions> <code>
# or
ssh -t bashform.me c <num_of_questions> <code>
```

Replace `<num_of_questions>` with the number of questions you want in the form, and `<code>` with the unique code for your form.

<p align="center">
    <img height="600" src="/images/create.gif">
</p>

### Get Forms and Responses

To get a list of forms and responses, use the following command:

```bash
ssh -t bashform.me forms
```

This will display a list of forms and by selecting a form, you can view the responses.

<p align="center">
    <img height="600" src="/images/forms.gif">
</p>

## Example

### Creating a Form

```bash
ssh -t bashform.me create 2 myform
```

This creates a form with 2 questions and the code `myform`.

### Filling Out a Form

```bash
ssh -t bashform.me form myform
```

This allows you to respond to the form with the code `myform`.

### Getting Form Responses

```bash
ssh -t bashform.me forms
```

## Try It Out

You can try out Bashform by using the following commands:

```bash
ssh -t bashform.me f devmegablaster
```

## Contributing

Contributions are welcome! If you'd like to improve Bashform, follow these steps:

1. Fork the repository.
2. Create a new branch: `git checkout -b feature-branch`.
3. Make your changes and commit: `git commit -m 'Add new feature'`.
4. Push to the branch: `git push origin feature-branch`.
5. Submit a pull request.
