import './App.css'

import {Button, Link, Navbar, NavbarBrand, NavbarContent, NavbarItem} from "@nextui-org/react";
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import {faPlus} from "@fortawesome/free-solid-svg-icons";

function App() {
  return (
    <>
        <Navbar>
            <NavbarBrand>
                <p className="font-bold text-inherit">Jobs Dashboard</p>
            </NavbarBrand>

            <NavbarContent justify="end">
                <NavbarItem>
                    <Button as={Link} color="primary" variant="flat" startContent={<FontAwesomeIcon icon={faPlus} />}>
                        Add Job
                    </Button>
                </NavbarItem>
            </NavbarContent>
        </Navbar>
    </>
  )
}

export default App
