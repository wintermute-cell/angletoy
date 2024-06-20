#!/usr/bin/python3

import argparse
from svgelements import SVG, Rect, Matrix


def extract_rect_details(rect_element):
    position = rect_element.x, rect_element.y
    size = rect_element.width, rect_element.height
    rotation = 0.0

    if rect_element.transform:
        m = Matrix(rect_element.transform)
        rotation = m.rotation

    origin = rect_element.x + (rect_element.width / 2), rect_element.y + (rect_element.height / 2)
    return position, size, rotation, origin


def main(svg_path):
    svg = SVG.parse(svg_path)
    for element in svg.elements():
        if isinstance(element, Rect):
            position, size, rotation, origin = extract_rect_details(element)
            print(element)
            print(f"position: {position[0]}, {position[1]} ; size: {size[0]}, {size[1]} ; rotation: {rotation:.1f} ; origin: {origin[0]}, {origin[1]} ;")


if __name__ == "__main__":
    parser = argparse.ArgumentParser(description="Extract and print details of rect elements from an SVG file.")
    parser.add_argument("svg_path", type=str, help="Path to the SVG file.")
    
    args = parser.parse_args()
    main(args.svg_path)

